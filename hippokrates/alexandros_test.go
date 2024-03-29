package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	"strings"
)

func (l *OdysseiaFixture) theWordIsQueried(word string) error {
	response, err := l.client.Alexandros().Search(word, "", "", "", TraceId)
	if err != nil {
		return err
	}

	var meroi models.ExtendedResponse
	err = json.NewDecoder(response.Body).Decode(&meroi)

	l.ctx = context.WithValue(l.ctx, ResponseBody, meroi)

	return nil
}

func (l *OdysseiaFixture) theWordIsQueriedWithAndNotFound(word string) error {
	_, err := l.client.Alexandros().Search(word, "", "exact", "", TraceId)
	if err != nil {
		l.ctx = context.WithValue(l.ctx, ErrorBody, err.Error())
	}

	return nil
}

func (l *OdysseiaFixture) theWordIsQueriedUsingAndAnd(word, mode, language, expand string) error {
	response, err := l.client.Alexandros().Search(word, language, mode, expand, TraceId)
	if err != nil {
		return err
	}

	var meroi models.ExtendedResponse
	err = json.NewDecoder(response.Body).Decode(&meroi)

	l.ctx = context.WithValue(l.ctx, ResponseBody, meroi)

	return nil
}

func (l *OdysseiaFixture) thePartialIsQueried(partial string) error {
	response, err := l.client.Alexandros().Search(partial, "", "partial", "", TraceId)
	if err != nil {
		return err
	}

	var meroi models.ExtendedResponse
	err = json.NewDecoder(response.Body).Decode(&meroi)

	l.ctx = context.WithValue(l.ctx, ResponseBody, meroi)

	return nil
}

func (l *OdysseiaFixture) theWordIsStrippedOfAccents(word string) error {
	strippedWord := transform.RemoveAccents(word)

	response, err := l.client.Alexandros().Search(strippedWord, "", "", "", TraceId)
	if err != nil {
		return err
	}

	var meroi models.ExtendedResponse
	err = json.NewDecoder(response.Body).Decode(&meroi)

	l.ctx = context.WithValue(l.ctx, ResponseBody, meroi)

	return nil
}

func (l *OdysseiaFixture) theWordShouldBeIncludedInTheResponse(searchTerm string) error {
	words := l.ctx.Value(ResponseBody).(models.ExtendedResponse)

	found := false

	for _, hit := range words.Hits {
		if strings.Contains(hit.Hit.Greek, searchTerm) {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("could not find searchterm %v in slice", searchTerm)
	}
	return nil
}

func (l *OdysseiaFixture) aGreekTranslationShouldBeIncludedInTheResponse() error {
	words := l.ctx.Value(ResponseBody).(models.ExtendedResponse)

	included := true

	for _, hit := range words.Hits {
		if hit.Hit.Greek == "" {
			included = false
		}
	}

	if !included {
		return fmt.Errorf("could not find greek word")
	}
	return nil
}

func (l *OdysseiaFixture) anErrorContainingIsReturned(message string) error {
	errorText := l.ctx.Value(ErrorBody).(string)
	if !strings.Contains(errorText, message) {
		return fmt.Errorf("expected %v to contain %v", errorText, message)
	}

	return nil
}
