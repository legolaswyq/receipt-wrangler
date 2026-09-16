# openapi.model.UpsertDashboardCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Dashboard name | 
**groupId** | **String** | Group foreign key | 
**widgets** | [**BuiltList&lt;UpsertWidgetCommand&gt;**](UpsertWidgetCommand.md) | Widgets associated to dashboard | [optional] 
**period** | **String** | Dashboard-level date range preset (THIS_MONTH, LAST_MONTH, LAST_3_MONTHS, THIS_YEAR, ALL_TIME, CUSTOM) | [optional] 
**periodStartDate** | **String** | ISO start date, used only when period is CUSTOM | [optional] 
**periodEndDate** | **String** | ISO end date, used only when period is CUSTOM | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


