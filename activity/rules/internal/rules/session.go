package rules

import (
	"fmt"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

// CreateRuleSession - creates a rule session
func (this *RuleEngine) CreateRuleSessionThenStart(tupleDescriptor string) error {

	log.Info(fmt.Sprintf("[RuleEngine.CreateRuleSessionThenStart] entering ... \n"))

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return err
	}
	this.ruleSession, err = ruleapi.GetOrCreateRuleSession("FlogoRulesSession")
	this.RegisterConditionsAndActions()
	this.ruleSession.Start(nil)

	log.Info(fmt.Sprintf("[RuleEngine.CreateRuleSessionThenStart] done ... \n"))

	return err
}

// CreateAndLoadRuleSession - creates a rule session and loads rules from config file
func (this *RuleEngine) CreateAndLoadRuleSessionThenStart(tupleDescriptor string, rulesDefs string) error {
	log.Info(fmt.Sprintf("Creating Rule Session and Rules \n"))

	model.RegisterTupleDescriptors(tupleDescriptor)

	this.RegisterConditionsAndActions()

	var err error
	this.ruleSession, err = ruleapi.GetOrCreateRuleSessionFromConfig("FlogoRulesSession", rulesDefs)
	this.ruleSession.Start(nil)
	return err
}
