# openapi.model.BankStatementImportTransaction

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**date** | **String** | Transaction date as YYYY-MM-DD | 
**name** | **String** | Receipt name | 
**amount** | **String** | Transaction amount, always positive | 
**isIncome** | **bool** | True for a deposit, false for a withdrawal | 
**categoryId** | **int** | Category to assign. Required, and must be an income category, when isIncome is true | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


