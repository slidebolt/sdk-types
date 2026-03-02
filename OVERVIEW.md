### `sdk-types` repository

#### Project Overview

This repository is a foundational component of the Slidebolt SDK. It provides the core data structures (types) that are used for communication and data exchange throughout the entire Slidebolt ecosystem.

#### Architecture

The `sdk-types` repository is a Go package that defines a set of shared data structures. It does not contain any application logic itself, but rather serves as a common language for all other Slidebolt components (such as the gateway, plugins, and the runner).

The key principles of this package are:

-   **Shared Data Structures**: It defines all the essential types, including `Device`, `Entity`, `Command`, `Event`, and the JSON-RPC structures for communication.
-   **Data Transfer Objects (DTOs)**: The types in this package are designed as pure Data Transfer Objects. They are simple data containers with minimal to no behavior, ensuring a clean separation of data and logic. A test (`TestDTOPurity`) is included to enforce this principle.
-   **Domain Schema**: It provides a mechanism for defining and registering `DomainDescriptor`s, which act as schemas for different types of entities (e.g., `switch`, `light`, `sensor`).

#### Key Files

| File | Description |
| :--- | :--- |
| `go.mod` | Defines the Go module for the package. |
| `types.go` | The central file containing all the core data structure definitions for the Slidebolt ecosystem. |
| `registry.go` | Implements a simple, in-memory registry for `DomainDescriptor`s, allowing components to look up the capabilities of different domains. |
| `types_test.go`| Contains a "purity test" to ensure that the DTOs in this package remain simple, behavior-free data structures. |

#### Available Commands

This is a library package and is not intended to be run directly. It is imported and used by other Slidebolt components.
