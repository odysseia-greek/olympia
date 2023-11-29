package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	configs "github.com/odysseia-greek/olympia/parmenides/config"
	"strings"
	"sync"
)

type ParmenidesHandler struct {
	Config *configs.Config
}

func (p *ParmenidesHandler) DeleteIndexAtStartUp() error {
	deleted, err := p.Config.Elastic.Index().Delete(p.Config.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", p.Config.Index, deleted))
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
	policyCreated, err := p.Config.Elastic.Policy().CreateHotPolicy(p.Config.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", p.Config.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", p.Config.PolicyName))
	err := p.createPolicyAtStartup()
	if err != nil {
		return err

	}
	indexMapping := p.Config.Elastic.Builder().QuizIndex(p.Config.PolicyName)
	created, err := p.Config.Elastic.Index().Create(p.Config.Index, indexMapping)
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
			Channel: p.Config.Channel,
		}

		if alternateChannel {
			msg.Channel = p.Config.DutchChannel
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
			res, err := p.Config.Elastic.Document().Bulk(buf, p.Config.Index)
			if err != nil {
				logging.Error(err.Error())
				return err
			}

			p.Config.Created = p.Config.Created + len(res.Items)
		}
	}
	return nil
}

// EnqueueTask sends a task to the Eupalinos queue
func (p *ParmenidesHandler) EnqueueTask(ctx context.Context, message *pb.Epistello) error {
	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())

	_, err = p.Config.Eupalinos.EnqueueMessage(ctx, message)
	return err
}
