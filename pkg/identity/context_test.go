package identity

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	identity := Identity{
		Login: "test",
	}
	ctx := ContextWithIdentity(context.Background(), &identity)
	identityFromContext, ok := FromContext(ctx)
	if ok && identity.Login == identityFromContext.Login {
		return
	}

	t.Error("Identity from context did not match the one passed to it")
}
