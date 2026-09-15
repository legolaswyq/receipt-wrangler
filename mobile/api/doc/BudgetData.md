# openapi.model.BudgetData

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**month** | **String** | The month the budget data covers, formatted YYYY-MM | [optional] 
**income** | **double** | Total income for the month | [optional] 
**spent** | **double** | Total expense spend for the month | [optional] 
**net** | **double** | Income minus spent | [optional] 
**untracked** | **double** | Expense spend touching no budgeted category | [optional] 
**categories** | [**BuiltList&lt;BudgetCategory&gt;**](BudgetCategory.md) | Per-category budget vs. actual for the month | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


