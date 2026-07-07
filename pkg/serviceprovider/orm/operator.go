package orm

import (
	"fmt"
	"strings"
)

// Operator 操作符接口.
type Operator interface {
	Build(column string, value any) (string, []any, error)
}

var operators = map[string]Operator{
	"=":           binaryOperator("="),
	"eq":          binaryOperator("="),
	"!=":          binaryOperator("<>"),
	"<>":          binaryOperator("<>"),
	"ne":          binaryOperator("<>"),
	">":           binaryOperator(">"),
	">=":          binaryOperator(">="),
	"<":           binaryOperator("<"),
	"<=":          binaryOperator("<="),
	"like":        likeOperator{mode: likeBoth},
	"left like":   likeOperator{mode: likeLeft},
	"right like":  likeOperator{mode: likeRight},
	"in":          inOperator{not: false},
	"not in":      inOperator{not: true},
	"between":     betweenOperator{not: false},
	"not between": betweenOperator{not: true},
	"is null":     nullOperator{not: false},
	"is not null": nullOperator{not: true},
}

// RegisterOperator 注册自定义操作符
func RegisterOperator(name string, operator Operator) {
	name = strings.ToLower(strings.TrimSpace(name))
	operators[name] = operator
}

// BuildOperator 构建SQL
func BuildOperator(column, operator string, value any) (string, []any, error) {
	operator = strings.ToLower(strings.TrimSpace(operator))

	op, ok := operators[operator]
	if !ok {
		return "", nil, fmt.Errorf("不支持操作符[%s]", operator)
	}

	return op.Build(column, value)
}

type binaryOperator string

func (o binaryOperator) Build(column string, value any) (string, []any, error) {
	return fmt.Sprintf("%s %s ?", column, string(o)), []any{value}, nil
}

const (
	likeLeft = iota
	likeRight
	likeBoth
)

type likeOperator struct {
	mode int
}

func (o likeOperator) Build(column string, value any) (string, []any, error) {
	text := fmt.Sprint(value)

	switch o.mode {
	case likeLeft:
		text += "%"
	case likeRight:
		text = "%" + text
	default:
		text = "%" + text + "%"
	}

	return fmt.Sprintf("%s LIKE ?", column), []any{text}, nil
}

type inOperator struct {
	not bool
}

func (o inOperator) Build(column string, value any) (string, []any, error) {
	values, ok := value.([]any)
	if !ok {
		return "", nil, fmt.Errorf("IN参数必须为数组")
	}

	if len(values) == 0 {
		return "", nil, fmt.Errorf("IN参数不能为空")
	}

	if o.not {
		return fmt.Sprintf("%s NOT IN ?", column), []any{values}, nil
	}

	return fmt.Sprintf("%s IN ?", column), []any{values}, nil
}

type betweenOperator struct {
	not bool
}

func (o betweenOperator) Build(column string, value any) (string, []any, error) {
	values, ok := value.([]any)
	if !ok {
		return "", nil, fmt.Errorf("BETWEEN参数必须为数组")
	}

	if len(values) != 2 {
		return "", nil, fmt.Errorf("BETWEEN参数必须包含两个值")
	}

	sql := "%s BETWEEN ? AND ?"

	if o.not {
		sql = "%s NOT BETWEEN ? AND ?"
	}

	return fmt.Sprintf(sql, column), []any{
		values[0],
		values[1],
	}, nil
}

type nullOperator struct {
	not bool
}

func (o nullOperator) Build(column string, value any) (string, []any, error) {
	if o.not {
		return fmt.Sprintf("%s IS NOT NULL", column), nil, nil
	}

	return fmt.Sprintf("%s IS NULL", column), nil, nil
}
