package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatStages: []testerutils.Stage{},
	ExecutableFileName: "your_grep.sh",
	Stages: []testerutils.Stage{
		{
			Slug:                    "init",
			Title:                   "Print number of tables",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
	}
}
