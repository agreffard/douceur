package css

import "fmt"

const (
	IDENT_SPACES = 2
)

// Rule kinds
const (
	QUALIFIED_RULE RuleKind = iota
	AT_RULE
)

// At Rules than have Rules inside their block instead of Declarations
var atRulesWithRulesBlock = []string{
	"@document", "@font-feature-values", "@keyframes", "@media", "@supports",
}

type RuleKind int

type Rule struct {
	Kind         RuleKind
	Prelude      string
	Declarations []*Declaration

	// At Rule name (eg: "@media")
	Name string

	// At Rule embedded rules
	Rules []*Rule

	// Current rule embedding level
	EmbedLevel int
}

func NewRule(kind RuleKind) *Rule {
	return &Rule{
		Kind: kind,
	}
}

// Returns string representation of rule kind
func (kind RuleKind) String() string {
	switch kind {
	case QUALIFIED_RULE:
		return "Qualified Rule"
	case AT_RULE:
		return "At Rule"
	default:
		return "WAT"
	}
}

// Returns true if this rule embeds another rules
func (rule *Rule) EmbedsRules() bool {
	if rule.Kind == AT_RULE {
		for _, atRuleName := range atRulesWithRulesBlock {
			if rule.Name == atRuleName {
				return true
			}
		}
	}

	return false
}

// Returns true if both rules are equals
func (rule *Rule) Equal(other *Rule) bool {
	if (rule.Kind != other.Kind) ||
		(rule.Prelude != other.Prelude) ||
		(rule.Name != other.Name) {
		return false
	}

	if (len(rule.Declarations) != len(other.Declarations)) ||
		(len(rule.Rules) != len(other.Rules)) {
		return false
	}

	for i, decl := range rule.Declarations {
		if !decl.Equal(other.Declarations[i]) {
			return false
		}
	}

	for i, rule := range rule.Rules {
		if !rule.Equal(other.Rules[i]) {
			return false
		}
	}

	return true
}

// Returns a string representation of rules differences
func (rule *Rule) Diff(other *Rule) []string {
	result := []string{}

	if rule.Kind != other.Kind {
		result = append(result, fmt.Sprintf("Kind: %s | %s", rule.Kind.String(), other.Kind.String()))
	}

	if rule.Prelude != other.Prelude {
		result = append(result, fmt.Sprintf("Prelude: \"%s\" | \"%s\"", rule.Prelude, other.Prelude))
	}

	if rule.Name != other.Name {
		result = append(result, fmt.Sprintf("Name: \"%s\" | \"%s\"", rule.Name, other.Name))
	}

	if len(rule.Declarations) != len(other.Declarations) {
		result = append(result, fmt.Sprintf("Declarations Nb: %d | %d", len(rule.Declarations), len(other.Declarations)))
	}

	if len(rule.Rules) != len(other.Rules) {
		result = append(result, fmt.Sprintf("Rules Nb: %d | %d", len(rule.Rules), len(other.Rules)))
	}

	for i, decl := range rule.Declarations {
		if !decl.Equal(other.Declarations[i]) {
			result = append(result, fmt.Sprintf("Declaration: \"%s\" | \"%s\"", decl.String(), other.Declarations[i].String()))
		}
	}

	for i, rule := range rule.Rules {
		if !rule.Equal(other.Rules[i]) {
			result = append(result, fmt.Sprintf("Rule: \"%s\" | \"%s\"", rule.String(), other.Rules[i].String()))
		}
	}

	return result
}

// Returns the string representation of a rule
func (rule *Rule) String() string {
	result := ""

	// result += fmt.Sprintf("[%s] ", rule.Kind.String())

	if rule.Kind == AT_RULE {
		result += fmt.Sprintf("%s", rule.Name)
	}

	if rule.Prelude != "" {
		if result != "" {
			result += " "
		}
		result += fmt.Sprintf("%s", rule.Prelude)
	}

	if (len(rule.Declarations) == 0) && (len(rule.Rules) == 0) {
		result += ";"
	} else {

		result += " {\n"

		if rule.EmbedsRules() {
			for _, subRule := range rule.Rules {
				result += fmt.Sprintf("%s%s\n", rule.indent(), subRule.String())
			}
		} else {
			for _, decl := range rule.Declarations {
				result += fmt.Sprintf("%s%s\n", rule.indent(), decl.String())
			}
		}

		result += fmt.Sprintf("%s}", rule.indentEndBlock())
	}

	return result
}

// Returns identation spaces for declarations and rules
func (rule *Rule) indent() string {
	result := ""

	for i := 0; i < ((rule.EmbedLevel + 1) * IDENT_SPACES); i++ {
		result += " "
	}

	return result
}

// Returns identation spaces for end of block character
func (rule *Rule) indentEndBlock() string {
	result := ""

	for i := 0; i < (rule.EmbedLevel * IDENT_SPACES); i++ {
		result += " "
	}

	return result
}