# openapi.model.UpsertReceiptProcessingSettingsCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Name of the settings | 
**description** | **String** | Description of the settings | [optional] 
**aiType** | [**AiType**](AiType.md) |  | 
**url** | **String** | URL for custom endpoints | [optional] 
**key** | **String** | Key for endpoints that require authentication | [optional] 
**model** | **String** | LLM model | [optional] 
**isVisionModel** | **bool** | Is vision model | [optional] 
**enforceJsonResponseFormat** | **bool** | Enforce JSON response format on the LLM provider. Disable if the provider does not support this flag. | [optional] 
**ocrEngine** | [**OcrEngine**](OcrEngine.md) |  | 
**ocrEngineUrl** | **String** | URL for the OCR engine's endpoint (used when OcrEngine is CUSTOM, e.g. a self-hosted Ollama vision model) | [optional] 
**ocrEngineModel** | **String** | Model for the OCR engine (used when OcrEngine is CUSTOM, e.g. a self-hosted Ollama vision model) | [optional] 
**promptId** | **int** | Prompt foreign key | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


