# openapi.api.ImportApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**importBankStatement**](ImportApi.md#importbankstatement) | **POST** /import/bankStatement | Import bank statement transactions
[**importConfigJson**](ImportApi.md#importconfigjson) | **POST** /import/importConfigJson | Import config json
[**previewBankStatement**](ImportApi.md#previewbankstatement) | **POST** /import/bankStatement/preview | Preview a bank statement


# **importBankStatement**
> BankStatementImportResult importBankStatement(bankStatementImportCommand)

Import bank statement transactions

Creates one resolved receipt per reviewed bank statement transaction

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getImportApi();
final BankStatementImportCommand bankStatementImportCommand = ; // BankStatementImportCommand | Reviewed transactions to import

try {
    final response = api.importBankStatement(bankStatementImportCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ImportApi->importBankStatement: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bankStatementImportCommand** | [**BankStatementImportCommand**](BankStatementImportCommand.md)| Reviewed transactions to import | 

### Return type

[**BankStatementImportResult**](BankStatementImportResult.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **importConfigJson**
> importConfigJson(file)

Import config json

This will import a config json

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getImportApi();
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | Files to quick scan

try {
    api.importConfigJson(file);
} catch on DioException (e) {
    print('Exception when calling ImportApi->importConfigJson: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | **MultipartFile**| Files to quick scan | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **previewBankStatement**
> BankStatementPreview previewBankStatement(file, groupId)

Preview a bank statement

Parses a bank statement PDF and returns its transactions for review. Nothing is created.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getImportApi();
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | Bank statement PDF
final int groupId = 56; // int | Group the transactions will be imported into

try {
    final response = api.previewBankStatement(file, groupId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ImportApi->previewBankStatement: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | **MultipartFile**| Bank statement PDF | 
 **groupId** | **int**| Group the transactions will be imported into | 

### Return type

[**BankStatementPreview**](BankStatementPreview.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

