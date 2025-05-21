from enum import Enum


class Model(str, Enum):
    GEMINI_2_5_FLASH_PREVIEW_04_17 = "gemini-2.5-flash-preview-04-17"
    GEMINI_2_5_PRO_EXPERIMENTAL_03_25 = "gemini-2.5-pro-exp-03-25"
    GEMINI_2_0_FLASH = "gemini-2.0-flash"
    GEMINI_2_0_FLASH_LITE = "gemini-2.0-flash-lite"

    def __str__(self):
        return self.value
