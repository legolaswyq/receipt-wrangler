# openapi.model.Dashboard

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdAt** | **String** |  | [optional] 
**createdBy** | **int** |  | [optional] 
**id** | **int** |  | 
**name** | **String** | Dashboard name | 
**groupId** | **int** | Group foreign key | [optional] 
**userId** | **int** | User foreign key | 
**updatedAt** | **String** |  | [optional] 
**widgets** | [**BuiltList&lt;Widget&gt;**](Widget.md) | Widgets associated to dashboard | [optional] 
**period** | **String** | Dashboard-level date range preset (THIS_MONTH, LAST_MONTH, LAST_3_MONTHS, THIS_YEAR, ALL_TIME, CUSTOM). Empty defaults to THIS_MONTH. | [optional] 
**periodStartDate** | **String** | ISO start date, used only when period is CUSTOM | [optional] 
**periodEndDate** | **String** | ISO end date, used only when period is CUSTOM | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


