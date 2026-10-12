# openapi.model.BankStatementTransaction

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**date** | **String** | Transaction date as YYYY-MM-DD | 
**type** | **String** | Transaction type as printed on the statement | 
**details** | **String** | Merchant or counterparty as printed on the statement | 
**reference** | **String** | Extra reference line, such as a masked card number | 
**name** | **String** | Suggested receipt name | 
**amount** | **String** | Transaction amount, always positive | 
**balance** | **String** | Running balance after the transaction | 
**isIncome** | **bool** | True for a deposit, false for a withdrawal | 
**selected** | **bool** | Whether the row is suggested for import by default | 
**duplicate** | **bool** | Whether a matching receipt already exists in the group | 
**suggestedCategoryId** | **int** | Suggested category, when one could be worked out | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


