# openapi.model.BankStatementPreview

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**bank** | **String** | Bank the statement was recognised as | 
**accountName** | **String** | Account name | 
**accountNumber** | **String** | Account number | 
**periodStart** | **String** | Statement period start as YYYY-MM-DD | 
**periodEnd** | **String** | Statement period end as YYYY-MM-DD | 
**transactions** | [**BuiltList&lt;BankStatementTransaction&gt;**](BankStatementTransaction.md) | Parsed transactions, newest first | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


