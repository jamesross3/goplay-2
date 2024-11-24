package celplay

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/require"
)

func TestCEL(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
		cel.Variable("group", cel.StringType),
	)
	require.NoError(t, err)
	require.NotEmpty(t, env)
	ast, issues := env.Compile(`name.startsWith("/groups/"+group)`)
	require.Nil(t, issues)
	require.NotEmpty(t, ast)

	program, err := env.Program(ast)
	require.NoError(t, err)
	refVal, evalDetails, err := program.Eval(map[string]any{
		"name":  "/groups/foo",
		"group": "foo",
	})

	require.NoError(t, err)
	require.NotEmpty(t, refVal)
	require.Nil(t, evalDetails)
	t.Logf("Ref val: %v\n", refVal.Value())
}
