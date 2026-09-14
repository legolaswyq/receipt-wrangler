package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

func setupBudgetsTest() {
	repositories.CreateTestGroupWithUsers()
	repositories.CreateTestCategories()
}

func tearDownBudgetsTest() {
	repositories.TruncateTestDb()
}

func newBudgetTestRequest(method, target string, body string, groupId string, categoryId string) (*httptest.ResponseRecorder, *http.Request) {
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, target, reader)

	rctx := chi.NewRouteContext()
	if groupId != "" {
		rctx.URLParams.Add("groupId", groupId)
	}
	if categoryId != "" {
		rctx.URLParams.Add("categoryId", categoryId)
	}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1},
	})
	r = r.WithContext(newContext)

	return w, r
}

func TestBudget_GetBudgetData_Success(t *testing.T) {
	defer tearDownBudgetsTest()
	setupBudgetsTest()

	grantAllGroupPerms(t, 1, 1)

	w, r := newBudgetTestRequest("POST", "/api/budget/1", "", "1", "")

	GetBudgetData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.BudgetData
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if result.Categories == nil {
		utils.PrintTestError(t, "categories is nil", "expected non-nil categories slice")
	}
}

func TestBudget_UpsertBudget_InvalidAmount(t *testing.T) {
	defer tearDownBudgetsTest()
	setupBudgetsTest()

	grantAllGroupPerms(t, 1, 1)

	requestBody := commands.UpsertCategoryBudgetCommand{
		CategoryId: 1,
		Amount:     decimal.Zero,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	w, r := newBudgetTestRequest("PUT", "/api/budget/1", string(body), "1", "")

	UpsertBudget(w, r)

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestBudget_UpsertBudget_Success(t *testing.T) {
	defer tearDownBudgetsTest()
	setupBudgetsTest()

	grantAllGroupPerms(t, 1, 1)

	requestBody := commands.UpsertCategoryBudgetCommand{
		CategoryId: 1,
		Amount:     decimal.NewFromFloat(100.00),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	w, r := newBudgetTestRequest("PUT", "/api/budget/1", string(body), "1", "")

	UpsertBudget(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}
