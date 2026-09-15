# openapi.api.BudgetApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**deleteBudget**](BudgetApi.md#deletebudget) | **DELETE** /budget/{groupId}/{categoryId} | Delete category budget
[**getBudgetData**](BudgetApi.md#getbudgetdata) | **POST** /budget/{groupId} | Get budget data
[**upsertBudget**](BudgetApi.md#upsertbudget) | **PUT** /budget/{groupId} | Upsert category budget


# **deleteBudget**
> deleteBudget(groupId, categoryId)

Delete category budget

This will delete a category's monthly budget target for a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getBudgetApi();
final int groupId = 56; // int | Id of group the budget belongs to
final int categoryId = 56; // int | Id of the category whose budget should be deleted

try {
    api.deleteBudget(groupId, categoryId);
} catch on DioException (e) {
    print('Exception when calling BudgetApi->deleteBudget: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of group the budget belongs to | 
 **categoryId** | **int**| Id of the category whose budget should be deleted | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getBudgetData**
> BudgetData getBudgetData(groupId)

Get budget data

This will get the current month's budget data for a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getBudgetApi();
final int groupId = 56; // int | Id of group to get or update budget data for

try {
    final response = api.getBudgetData(groupId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling BudgetApi->getBudgetData: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of group to get or update budget data for | 

### Return type

[**BudgetData**](BudgetData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **upsertBudget**
> CategoryBudget upsertBudget(groupId, upsertCategoryBudgetCommand)

Upsert category budget

This will create or update a category's monthly budget target for a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getBudgetApi();
final int groupId = 56; // int | Id of group to get or update budget data for
final UpsertCategoryBudgetCommand upsertCategoryBudgetCommand = ; // UpsertCategoryBudgetCommand | Category budget data to upsert

try {
    final response = api.upsertBudget(groupId, upsertCategoryBudgetCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling BudgetApi->upsertBudget: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of group to get or update budget data for | 
 **upsertCategoryBudgetCommand** | [**UpsertCategoryBudgetCommand**](UpsertCategoryBudgetCommand.md)| Category budget data to upsert | 

### Return type

[**CategoryBudget**](CategoryBudget.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

