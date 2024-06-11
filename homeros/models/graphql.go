package models

import (
	"github.com/odysseia-greek/agora/plato/models"
)

type Health struct {
	Overall    bool          `json:"overallHealth"`
	Herodotos  models.Health `json:"herodotos"`
	Sokrates   models.Health `json:"sokrates"`
	Dionysios  models.Health `json:"dionysios"`
	Alexandros models.Health `json:"alexandros"`
}
