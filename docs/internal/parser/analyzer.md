# Internal Parser

This module handles the AST (Abstract Syntax Tree) analysis of Go source files.

## Functions

### AnalyzeFile
Extracts specific functions and structs that intersect with changed lines provided by the Git diff.

### CalculateCoverage
Traverses the project to count exported symbols and verify if they have associated documentation, either as comments or in the `docs/` directory.
