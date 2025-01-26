package ollama_test

import (
	"testing"

	"github.com/k0kubun/pp"

	"github.com/lakrizz/prollama/pkg/ollama"
)

//
// func TestPatchComment(t *testing.T) {
// 	raw, err := os.ReadFile("test_patch")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	pp.Println(ollama.FindLineNumber("lg.Error(\"failde creating schema resources: %v\", \"error\", err)", string(raw)))
// }

func TestFindLineNumber(t *testing.T) {
	needle := "+\tMailtrap      *mailer.Mailtrap"
	haystack := "index 1836d151..9a6a12c2 100644\n--- a/backend/api/router.go\n+++ b/backend/api/router.go\n@@ -24,6 +24,7 @@ import (\n \t\"hooks.im/api/db/webhook\"\n \t\"hooks.im/api/internal/hash\"\n \t\"hooks.im/api/internal/jwt\"\n+\t\"hooks.im/api/internal/mailer\"\n \t\"hooks.im/api/internal/paypal\"\n \t\"hooks.im/api/internal/ptr\"\n \t\"hooks.im/api/internal/subscription\"\n@@ -43,6 +44,30 @@ type Router struct {\n \tTargets       *targets.Targets\n \tPayPalService *paypal.Service\n \tJwtService    *jwt.JwtService\n+\tMailtrap      *mailer.Mailtrap\n+}\n+\n+// Contact form submission\n+// (POST /landing/contact/)\n+func (ro *Router) SubmitContactForm(ctx echo.Context) error {\n+\treq := gen.SubmitContactFormJSONRequestBody{}\n+\terr := ctx.Bind(&req)\n+\tif err != nil {\n+\t\treturn ctx.JSON(http.StatusBadRequest, &gen.ErrorResponse{\n+\t\t\tMessage: \"error binding contact form data\",\n+\t\t\tDetails: ptr.String(err.Error()),\n+\t\t})\n+\t}\n+\n+\terr = ro.Mailtrap.SendMail(mailer.TEMPLATE_CONTACTFORM, \"contact@piles.dev\", map[string]any{\"name\": req.Name, \"email\": req.Email, \"body\": req.Message})\n+\tif err != nil {\n+\t\treturn ctx.JSON(http.StatusBadRequest, &gen.ErrorResponse{\n+\t\t\tMessage: \"error sending mail\",\n+\t\t\tDetails: ptr.String(err.Error()),\n+\t\t})\n+\t}\n+\n+\treturn ctx.JSON(http.StatusOK, nil)\n }\n \n // deletes a user with all associated resources"

	pp.Println(ollama.FindLineNumber(needle, haystack))
}
