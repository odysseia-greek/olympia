package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

const (
	defaultIndex string = "dictionary"
)

func CreateNewConfig(ctx context.Context) (*DiogenesHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer, err := aristophanes.NewClientTracer()
	if err != nil {
		logging.Error(err.Error())
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	ambassador := diplomat.NewClientAmbassador()
	ambassadorHealthy := ambassador.WaitForHealthyState()
	if !ambassadorHealthy {
		logging.Info("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceID := uuid.New().String()
	spanID := aristophanes.GenerateSpanID()
	combinedID := fmt.Sprintf("%s+%s+%d", traceID, spanID, 1)

	ambassadorCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	payload := &pbar.StartTraceRequest{
		Method:        "GetSecret",
		Url:           diplomat.DEFAULTADDRESS,
		Host:          "",
		RemoteAddress: "",
		Operation:     "/delphi_ptolemaios.Ptolemaios/GetSecret",
	}

	go func() {
		parabasis := &pbar.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pbar.ParabasisRequest_StartTrace{
				StartTrace: payload,
			},
		}
		if err := streamer.Send(parabasis); err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s", traceID, spanID))
	}()

	md := metadata.New(map[string]string{service.HeaderKey: combinedID})
	ambassadorCtx = metadata.NewOutgoingContext(context.Background(), md)
	vaultConfig, err := ambassador.GetSecret(ambassadorCtx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	go func() {
		parabasis := &pbar.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pbar.ParabasisRequest_CloseTrace{
				CloseTrace: &pbar.CloseTraceRequest{
					ResponseBody: fmt.Sprintf("user retrieved from vault: %s", vaultConfig.ElasticUsername),
				},
			},
		}

		err := streamer.Send(parabasis)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
	}()

	elasticService := aristoteles.ElasticService(tls)

	cfg := models.Config{
		Service:     elasticService,
		Username:    vaultConfig.ElasticUsername,
		Password:    vaultConfig.ElasticPassword,
		ElasticCERT: vaultConfig.ElasticCERT,
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = aristoteles.HealthCheck(elastic)
	if err != nil {
		return nil, err
	}

	englishToGreekDictionary, err := createEnglishToGreekDict()
	if err != nil {
		return nil, err
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)

	ctx, cancel := context.WithCancel(ctx)

	return &DiogenesHandler{
		Elastic:            elastic,
		Index:              index,
		Streamer:           streamer,
		Cancel:             cancel,
		EnglishToGreekDict: englishToGreekDictionary,
	}, nil
}

func createEnglishToGreekDict() (map[string]string, error) {
	const translationMap = `
        {
            "a": "α",  "a/": "ά",  "a\\": "ὰ", "a=": "ᾶ", "a(": "ἀ", "a)": "ἁ", "a|": "ᾳ", 
            "A": "Α",  "A/": "Ά",  "A\\": "Ὰ", "A=": "ᾶ", "A(": "Ἀ", "A)": "Ἁ", "A|": "ᾼ",
            "b": "β",  "B": "Β",
			"c": "κ",  "C": "Κ",
            "d": "δ",  "D": "Δ",
            "e": "ε",  "e/": "έ",  "e\\": "ὲ", "e(": "ἐ", "e)": "ἑ",
            "E": "Ε",  "E/": "Έ",  "E\\": "Ὲ", "E(": "Ἐ", "E)": "Ἑ",
            "f": "φ",  "F": "Φ",
            "g": "γ",  "G": "Γ",
            "h": "η",  "h/": "ή",  "h\\": "ὴ", "h=": "ῆ", "h(": "ἠ", "h)": "ἡ", "h|": "ῃ",
            "H": "Η",  "H/": "Ή",  "H\\": "Ὴ", "H=": "ῆ", "H(": "Ἠ", "H)": "Ἡ", "H|": "ῌ",
            "i": "ι",  "i/": "ί",  "i\\": "ὶ", "i=": "ῖ", "i(": "ἰ", "i)": "ἱ",
            "I": "Ι",  "I/": "Ί",  "I\\": "Ὶ", "I=": "ῖ", "I(": "Ἰ", "I)": "Ἱ",
            "j": "ξ",  "J": "Ξ",
            "k": "κ",  "K": "Κ",
            "l": "λ",  "L": "Λ",
            "m": "μ",  "M": "Μ",
            "n": "ν",  "N": "Ν",
            "o": "ο",  "o/": "ό",  "o\\": "ὸ", "o(": "ὀ", "o)": "ὁ",
            "O": "Ο",  "O/": "Ό",  "O\\": "Ὸ", "O(": "Ὀ", "O)": "Ὁ",
            "p": "π",  "P": "Π",
            "q": "θ",  "Q": "Θ",
            "r": "ρ",  "r(": "ῤ", "r)": "ῥ",
            "R": "Ρ",  "R(": "Ῥ",
            "s": "σ",  "s_end": "ς", "S": "Σ",
            "t": "τ",  "T": "Τ",
            "u": "υ",  "u/": "ύ",  "u\\": "ὺ", "u=": "ῦ", "u(": "ὐ", "u)": "ὑ", "u|": "ῡ",
            "U": "Υ",  "U/": "Ύ",  "U\\": "Ὺ", "U=": "ῦ", "U(": "Ὑ", "U)": "Ὑ", "U|": "Ῡ",
            "w": "ω",  "w/": "ώ",  "w\\": "ὼ", "w=": "ῶ", "w(": "ὠ", "w)": "ὡ", "w|": "ῳ",
            "W": "Ω",  "W/": "Ώ",  "W\\": "Ὼ", "W=": "ῶ", "W(": "Ὠ", "W)": "Ὡ", "W|": "ῼ",
            "x": "χ",  "X": "Χ",
            "y": "ψ",  "Y": "Ψ",
            "z": "ζ",  "Z": "Ζ"
        }`

	var dict map[string]string
	err := json.Unmarshal([]byte(translationMap), &dict)
	if err != nil {
		return nil, err
	}

	return dict, nil
}
