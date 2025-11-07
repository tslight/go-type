Internal Package Boundaries
===========================

This document explains the purpose of each internal package after refactor:

content/    Loading and management of embedded textual content (books, docs). No UI logic.
state/      Persistence layer for progress and session statistics (state file handling).
selection/  Interactive selection flow plus lightweight state provider implementation.
menu/       Bubble Tea menu UI for choosing content and viewing stats.
model/      Typing session state machine and rendering logic.
runner/     Application harness wiring flags, selection, and model execution.
utils/      Pure calculation helpers (WPM, accuracy, error counts).

Design Notes:
-------------
- Provider implementation is located in selection to avoid an import cycle between content <-> state.
- model depends only on a narrow SessionState interface for persistence, supporting test doubles.
- cmd/ binaries import runner + selection + content to compose the application.
- Deprecated statestore/ retained temporarily only during migration; slated for removal.

Extensibility:
--------------
- To add a new content source variant: create a new embedded FS and instantiate ContentManager.
- To replace persistence: implement SessionState and inject into model.NewModel.
- To add analytics: extend state.ContentState / SessionResult (since internal, refactors are safe).
