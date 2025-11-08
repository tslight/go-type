# Typing Statistics & WPM Calculation

This document explains how typing statistics are computed, the difference between raw and effective character counts, session baselines, and how text progress is reported.

## Overview
The application tracks typing performance across sessions for each content item (book, doc, etc.). When you finish a typing session, summary metrics are recorded and aggregated.

You can also view aggregated statistics across all content items (global stats). This provides overall totals (sessions, time, characters) and weighted averages (WPM, accuracy) without relying on previously stored derived values.

## Key Concepts
- **Raw Characters**: Total runes you typed during the current session. This includes spaces, newlines, and all visible characters. It excludes any pre-filled characters loaded from previous progress.
- **Effective Characters**: Filtered character count that excludes "excessive" whitespace (runs of >=3 spaces, duplicate newlines, etc.). This aligns with progress tracking.
- **Baseline**: When you resume a content item mid-way, previously completed text is pre-filled. We set baseline values (raw & effective) at session start so per-session metrics only count new typing.
- **Text Progress**: Absolute position in the source text (byte index +1) of the latest effective character typed, reported as `Text Progress: current/total`.

## WPM Calculation
We use the common formula:

```
WPM = (characters_typed / 5) / minutes_elapsed
```

Details:
- Characters counted are UTF-8 runes (so multi-byte characters don’t inflate counts).
- Only characters typed **this session** (post-baseline) contribute to WPM.
- For very short sessions: if you typed any characters but elapsed time is under 1 second, we compute WPM with a minimum of 1 second to avoid a misleading 0.00 value.

## Accuracy & Errors
- **Accuracy (%)**: `(correct_characters / total_text_length) * 100`. It compares your input against the source text up to the shorter length. Characters beyond the text count as errors.
- **Errors**: Mismatched characters plus extra or missing characters relative to the source length.

## Session Recording
Each session stores:
- Timestamp
- WPM
- Accuracy
- Errors
- Raw characters typed this session
- Effective characters typed this session
- Duration (seconds)

Stored JSON fields in each session:
```json
{
  "timestamp": "2025-11-08T12:34:56Z",
  "wpm": 42.3,
  "accuracy": 97.5,
  "errors": 12,
  "characters_typed_raw": 210,
  "characters_typed_effective": 198,
  "duration_seconds": 300
}
```

## Aggregated Statistics (Per Content)
For each content item we maintain totals and averages:
- Sessions Completed
- Total Time
- Average WPM (time-weighted): computed from totals as `(sum(raw_chars)/5) / (sum(duration_seconds)/60)`
- Best WPM (highest recomputed per-session WPM)
- Average Accuracy
- Total Characters (raw)
- Total Characters (effective)
- Text Progress (current/total)

Raw and effective totals allow you to distinguish between literal keystrokes and meaningful progress.

## Global Statistics
Global statistics roll up every session across all content items:
- Sessions Completed (sum)
- Total Time (sum)
- Average WPM (time-weighted): `(sum(raw_chars)/5) / (sum(duration_seconds)/60)`
- Best WPM (max of recomputed per-session WPM values)
- Average Accuracy (mean of per-session accuracy values)
- Total Characters (raw) (sum)
- Total Characters (effective) (sum)
- Text Progress Typed (sum of per-content progress positions)
- Text Progress Total (sum of per-content total lengths)

Recomputation: For each stored session we recompute WPM from its raw character count and duration. For averages, we prefer time-weighted totals (see formulas above) rather than averaging per-session WPM values. We never use previously stored WPM values that might have been inflated before baseline logic was added.

Viewing: In the menu press `I` (capital i) for global stats, or `i` for the currently selected content item's stats.

## Examples
### Simple New Session
- Typed 250 characters over 2 minutes.
- Effective characters: 240 (10 characters were excessive whitespace).
- WPM = `(250 / 5) / 2 = 25 WPM`.

### Resumed Session (Baseline)
- You resume at character 10,000.
- You type 500 new characters in 1 minute.
- WPM uses only the 500 new characters: `(500 / 5) / 1 = 100 WPM`.
- Raw total in state aggregates +500; effective depends on filtering.

### Whitespace Heavy Session
- Raw: 100 characters
- Effective: 60 characters (40 were collapsed/excessive)
- WPM still uses raw session characters (standard practice), progress percent uses effective.

## Design Rationale
- Separating raw vs effective prevents inflated progress metrics while retaining realistic speed measurement.
- Baselines avoid counting historical progress in new session calculations, eliminating exaggerated WPM values when resuming.
- Using runes ensures international text produces accurate WPM.

## Future Enhancements
- Option to toggle WPM based on effective characters.
- Per-session diff view of errors.
- Median WPM addition alongside average.

---
If you spot a number that seems off, capture: total time, raw/effective counts, and the displayed WPM. With those we can recompute and verify consistency.

## Navigation Hint
While typing, press ESC to end the current session and return to the menu (results are recorded like Ctrl+Q/Ctrl+S but the program stays open).
