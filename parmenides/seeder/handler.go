package seeder

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
	"sync"
)

type ParmenidesHandler struct {
	Index        string
	Created      int
	Elastic      elastic.Client
	Eupalinos    EupalinosClient
	Channel      string
	DutchChannel string
	ExitCode     string
	PolicyName   string
	Ambassador   *ptolemaios.ClientAmbassador
}

func (p *ParmenidesHandler) DeleteIndexAtStartUp() error {
	deleted, err := p.Elastic.Index().Delete(p.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", p.Index, deleted))
	if err != nil {
		if deleted {
			return nil
		}
		if strings.Contains(err.Error(), "index_not_found_exception") {
			logging.Error(err.Error())
			return nil
		}

		return err
	}

	return nil
}

func (p *ParmenidesHandler) createPolicyAtStartup() error {
	policyCreated, err := p.Elastic.Policy().CreateHotPolicy(p.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", p.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", p.PolicyName))
	err := p.createPolicyAtStartup()
	if err != nil {
		return err

	}
	indexMapping := p.Elastic.Builder().QuizIndex(p.PolicyName)
	created, err := p.Elastic.Index().Create(p.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) Add(logoi models.Logos, wg *sync.WaitGroup, method, category string) error {
	defer wg.Done()
	var buf bytes.Buffer

	var currBatch int

	for _, word := range logoi.Logos {
		currBatch++

		meros := models.Meros{
			Greek:      word.Greek,
			English:    word.Translation,
			LinkedWord: "",
			Original:   word.Greek,
		}

		alternateChannel := false
		if method == "mouseion" {
			meros.Dutch = word.Translation
			meros.English = ""
			alternateChannel = true
		}

		jsonsifiedMeros, _ := meros.Marshal()
		msg := &pb.Epistello{
			Data:    string(jsonsifiedMeros), // Your JSON data here
			Channel: p.Channel,
		}

		if alternateChannel {
			msg.Channel = p.DutchChannel
		}

		err := p.EnqueueTask(context.Background(), msg)
		if err != nil {
			logging.Error(err.Error())
		}

		word.Category = category
		word.Method = method

		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedWord, _ := word.Marshal()
		jsonifiedWord = append(jsonifiedWord, "\n"...)
		buf.Grow(len(meta) + len(jsonifiedWord))
		buf.Write(meta)
		buf.Write(jsonifiedWord)

		if currBatch == len(logoi.Logos) {
			res, err := p.Elastic.Document().Bulk(buf, p.Index)
			if err != nil {
				logging.Error(err.Error())
				return err
			}

			p.Created = p.Created + len(res.Items)
		}
	}
	return nil
}

// EnqueueTask sends a task to the Eupalinos queue
func (p *ParmenidesHandler) EnqueueTask(ctx context.Context, message *pb.Epistello) error {
	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())

	_, err = p.Eupalinos.EnqueueMessage(ctx, message)
	return err
}
