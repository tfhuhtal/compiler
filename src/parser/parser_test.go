package parser_test

//import (
//"compiler/src/ast"
//"compiler/src/parser"
//"compiler/src/tokenizer"
//"testing"
//)

//func TestParseIntLiteral(t *testing.T) {
//toks := []tokenizer.Token{
//{Type: "int_literal", Text: "123", Location: "loc1"},
//}
//lit, err := parser.Parse(toks)
//if err != nil {
//t.Fatalf("unexpected error: %v", err)
//}
//intVal, ok := lit.(ast.Literal)
//if !ok {
//t.Fatalf("expected ast.Literal, got %T", lit)
//}
//if intVal.Value != 123 {
//t.Fatalf("expected 123, got %d", intVal.Value)
//}
//}

//func TestParseExpressionAdd(t *testing.T) {
//toks := []tokenizer.Token{
//{Type: "int_literal", Text: "1", Location: "loc1"},
//{Type: "operator", Text: "+", Location: "loc2"},
//{Type: "int_literal", Text: "2", Location: "loc3"},
//}
//expr, err := parser.Parse(toks)
//if err != nil {
//t.Fatalf("unexpected error: %v", err)
//}
//binOp, ok := expr.(ast.BinaryOp)
//if !ok {
//t.Fatalf("expected ast.BinaryOp, got %T", expr)
//}
//if binOp.Left.Value != 1 || binOp.Right.Value != 2 || binOp.Op != "+" {
//t.Fatalf("expected 1 + 2, got %v %v %v", binOp.Left.Value, binOp.Op, binOp.Right.Value)
//}
//}

//func TestParseExpressionSub(t *testing.T) {
//toks := []tokenizer.Token{
//{Type: "int_literal", Text: "10", Location: "loc1"},
//{Type: "operator", Text: "-", Location: "loc2"},
//{Type: "int_literal", Text: "5", Location: "loc3"},
//}
//expr, err := parser.Parse(toks)
//if err != nil {
//t.Fatalf("unexpected error: %v", err)
//}
//binOp, ok := expr.(ast.BinaryOp)
//if !ok {
//t.Fatalf("expected ast.BinaryOp, got %T", expr)
//}
//if binOp.Left.Value != 10 || binOp.Right.Value != 5 || binOp.Op != "-" {
//t.Fatalf("expected 10 - 5, got %v %v %v", binOp.Left.Value, binOp.Op, binOp.Right.Value)
//}
//}

//func TestParseInvalidToken(t *testing.T) {
//toks := []tokenizer.Token{
//{Type: "identifier", Text: "abc", Location: "loc1"},
//}
//_, err := parser.Parse(toks)
//if err == nil {
//t.Fatal("expected error for invalid token, got nil")
//}
//}
