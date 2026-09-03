# openapi.model.ReceiptProcessingSettings

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**createdAt** | **String** |  | 
**createdBy** | **int** |  | [optional] [default to 0]
**createdByString** | **String** | Created by entity's name | [optional] [default to '']
**updatedAt** | **String** |  | [optional] [default to '']
**name** | **String** | Name of the settings | [optional] 
**description** | **String** | Description of the settings | [optional] 
**aiType** | [**AiType**](AiType.md) |  | [optional] 
**url** | **String** | URL for custom endpoints | [optional] 
**key** | **String** | Key for endpoints that require authentication | [optional] 
**model** | **String** | LLM model | [optional] 
**isVisionModel** | **bool** | Is vision model | [optional] 
**enforceJsonResponseFormat** | **bool** | Enforce JSON response format on the LLM provider. Disable if the provider does not support this flag. | [optional] 
**ocrEngine** | [**OcrEngine**](OcrEngine.md) |  | [optional] 
**ocrEngineUrl** | **String** | URL for the OCR engine's endpoint (used when OcrEngine is CUSTOM, e.g. a self-hosted Ollama vision model) | [optional] 
**ocrEngineModel** | **String** | Model for the OCR engine (used when OcrEngine is CUSTOM, e.g. a self-hosted Ollama vision model) | [optional] 
**prompt** | [**Prompt**](Prompt.md) |  | [optional] 
**promptId** | **int** | Prompt foreign key | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


