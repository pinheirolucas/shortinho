package database

import "errors"

// Engine is the enum type for databases
type Engine string

// Available options for Engine enum
const (
	EngineUnknown  Engine = ""
	EngineMongo    Engine = "mongo"
	EnginePostgres Engine = "postgres"
)

var validEngines = []Engine{EngineMongo, EnginePostgres}

// EngineFromString returns an engine from its string representation
func EngineFromString(engineStr string) (Engine, error) {
	engine := Engine(engineStr)
	var valid bool

	for _, validEngine := range validEngines {
		if validEngine == engine {
			valid = true
			break
		}
	}

	if valid {
		return engine, nil
	}

	return "", errors.New("invalid database engine")
}
