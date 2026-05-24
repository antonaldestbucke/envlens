// Package differ computes structured line-level diffs between reference
// and target .env key sets.
//
// Each diff line is annotated with a LineKind — context (shared), missing
// (absent from target), or extra (not in reference) — enabling downstream
// renderers and reporters to present precise, actionable output.
//
// Usage:
//
//	result := differ.Compute("production", refKeys, targetKeys)
//	for _, line := range result.Lines {
//		fmt.Println(line.Kind, line.Key)
//	}
package differ
