package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func bankStatementRequest(body *bytes.Buffer, contentType string, userId uint) *http.Request {
	r := httptest.NewRequest("POST", "/api/import/bankStatement", body)
	r.Header.Set("Content-Type", contentType)
	claims := &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}}
	return r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claims))
}

func bankStatementPreviewRequest(t *testing.T, userId uint, groupId string, file []byte) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if len(groupId) > 0 {
		_ = writer.WriteField("groupId", groupId)
	}
	if file != nil {
		part, err := writer.CreateFormFile("files", "statement.pdf")
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		_, _ = part.Write(file)
	}
	_ = writer.Close()

	return bankStatementRequest(body, writer.FormDataContentType(), userId)
}

func seedBankStatementGroup(t *testing.T) models.Group {
	t.Helper()
	group := models.Group{BaseModel: models.BaseModel{ID: 1}, Name: "statement-group"}
	if err := repositories.GetDB().Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	for _, id := range []uint{1, 2} {
		user := models.User{BaseModel: models.BaseModel{ID: id}, Username: fmt.Sprintf("statement-user-%d", id), Password: "p"}
		if err := repositories.GetDB().Create(&user).Error; err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}
	return group
}

func TestPreviewBankStatement_RejectsBadRequests(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedBankStatementGroup(t)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsCreate)

	tests := map[string]*http.Request{
		"missing group": bankStatementPreviewRequest(t, 1, "", []byte("x")),
		"missing file":  bankStatementPreviewRequest(t, 1, "1", nil),
		"empty file":    bankStatementPreviewRequest(t, 1, "1", []byte{}),
	}

	for name, request := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			PreviewBankStatement(w, request)
			if w.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Result().StatusCode)
			}
		})
	}
}

func TestPreviewBankStatement_RequiresReceiptCreatePermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedBankStatementGroup(t)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	w := httptest.NewRecorder()
	PreviewBankStatement(w, bankStatementPreviewRequest(t, 1, "1", []byte("x")))
	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Result().StatusCode)
	}
}

func TestPreviewBankStatement_RejectsUnsupportedFileWithMessage(t *testing.T) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext is not installed")
	}
	defer repositories.TruncateTestDb()
	seedBankStatementGroup(t)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsCreate)

	w := httptest.NewRecorder()
	PreviewBankStatement(w, bankStatementPreviewRequest(t, 1, "1", []byte("this is not a pdf")))
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
	if !strings.Contains(w.Body.String(), "not a supported bank statement") {
		t.Errorf("expected the unsupported-statement message, got %s", w.Body.String())
	}
}

func bankStatementImportRequest(t *testing.T, userId uint, payload any) *http.Request {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return bankStatementRequest(bytes.NewBuffer(body), "application/json", userId)
}

func TestImportBankStatement_CreatesReceipts(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedBankStatementGroup(t)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsCreate)
	db := repositories.GetDB()
	income := models.Category{Name: "Salary", IsIncome: true}
	if err := db.Create(&income).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}

	payload := map[string]any{
		"groupId":      1,
		"paidByUserId": 1,
		"transactions": []map[string]any{
			{"date": "2026-01-20", "name": "Corner Groc", "amount": "42.19", "isIncome": false},
			{"date": "2026-01-15", "name": "Employer Limited", "amount": "1500.00", "isIncome": true, "categoryId": income.ID},
		},
	}

	w := httptest.NewRecorder()
	ImportBankStatement(w, bankStatementImportRequest(t, 1, payload))
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Result().StatusCode, w.Body.String())
	}

	result := structs.BankStatementImportResult{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.CreatedCount != 2 {
		t.Errorf("expected createdCount 2, got %d", result.CreatedCount)
	}

	var count int64
	db.Model(&models.Receipt{}).Where("group_id = ?", 1).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 receipts, found %d", count)
	}
}

func TestImportBankStatement_RejectsInvalidAndUnauthorizedRequests(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedBankStatementGroup(t)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsCreate)
	grantGroupPerms(t, 2, 1, permissions.GroupReceiptsRead)
	db := repositories.GetDB()
	expense := models.Category{Name: "Groceries"}
	if err := db.Create(&expense).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}

	row := func(overrides map[string]any) map[string]any {
		base := map[string]any{"date": "2026-01-20", "name": "Corner Groc", "amount": "42.19", "isIncome": false}
		for key, value := range overrides {
			base[key] = value
		}
		return base
	}
	command := func(rows ...map[string]any) map[string]any {
		return map[string]any{"groupId": 1, "paidByUserId": 1, "transactions": rows}
	}

	tests := map[string]struct {
		userId  uint
		payload any
		expect  int
	}{
		"no transactions":                  {1, command(), http.StatusBadRequest},
		"missing name":                     {1, command(row(map[string]any{"name": " "})), http.StatusBadRequest},
		"bad date":                         {1, command(row(map[string]any{"date": "20/01/2026"})), http.StatusBadRequest},
		"zero amount":                      {1, command(row(map[string]any{"amount": "0"})), http.StatusBadRequest},
		"negative amount":                  {1, command(row(map[string]any{"amount": "-5"})), http.StatusBadRequest},
		"deposit without a category":       {1, command(row(map[string]any{"isIncome": true})), http.StatusBadRequest},
		"deposit with an expense category": {1, command(row(map[string]any{"isIncome": true, "categoryId": expense.ID})), http.StatusBadRequest},
		"missing group":                    {1, map[string]any{"paidByUserId": 1, "transactions": []map[string]any{row(nil)}}, http.StatusBadRequest},
		"missing payer":                    {1, map[string]any{"groupId": 1, "transactions": []map[string]any{row(nil)}}, http.StatusBadRequest},
		"malformed body":                   {1, "not an object", http.StatusBadRequest},
		"without create permission":        {2, command(row(nil)), http.StatusForbidden},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ImportBankStatement(w, bankStatementImportRequest(t, test.userId, test.payload))
			if w.Result().StatusCode != test.expect {
				t.Errorf("expected %d, got %d: %s", test.expect, w.Result().StatusCode, w.Body.String())
			}

			var count int64
			db.Model(&models.Receipt{}).Count(&count)
			if count != 0 {
				t.Errorf("expected no receipts to be created, found %d", count)
			}
		})
	}
}
