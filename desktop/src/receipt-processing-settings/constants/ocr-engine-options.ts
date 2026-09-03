import { FormOption } from "../../interfaces/form-option.interface";
import { OcrEngine } from "../../open-api";

export const ocrEngineOptions: FormOption[] = [
  {
    value: OcrEngine.Tesseract,
    displayValue: "Tesseract",
  },
  {
    value: OcrEngine.EasyOcr,
    displayValue: "EasyOCR",
  },
  {
    value: OcrEngine.Custom,
    displayValue: "Custom",
  }
];
