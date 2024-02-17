# Skeleton Code with Hexagonal Architecture

This repository contains a skeleton codebase with a Hexagonal Architecture approach for building various types of applications, including REST APIs, long-running process workers, serverless functions, etc. The purpose of this skeleton code is to provide a structured foundation, configuration setup, and common patterns for rapid development of new applications. By using this skeleton code, developers can streamline the creation process and focus more on solving business issues rather than setting up boilerplate code from scratch.

## Key Features
- **Hexagonal Architecture**: The codebase is organized following the principles of Hexagonal Architecture (also known as Ports and Adapters Architecture), which promotes separation of concerns and decoupling of application components.

> The implementation of Hexagonal Architecture in this repository following or heavily inspired by [Hex Monscape project](https://github.com/Haraj-backend/hex-monscape).

- **Programming Language Agnostic**: The skeleton code is designed to be language-agnostic, meaning it can be implemented in various programming languages. This allows developers to choose the language that best fits the requirements of their project.
- **Configurable Structure**: The codebase provides a configurable structure that can be adapted to different types of applications. It includes directories for domain logic, application services, adapters (e.g., HTTP, database), and configuration files.
- **Ready-to-Use Configuration**: Common configuration files, such as environment variables, logging setup, and dependency injection, are included to facilitate quick setup and deployment of applications.

## Usage
To use this skeleton code for your new application:

1. Clone this repository to your local machine.
2. Choose the programming language in which you want to implement your application. You can either use one of the existing language implementations provided in the repository or create a new implementation for your preferred language.
3. Customize the codebase according to your application requirements. Modify the domain logic, implement application services, and configure adapters as needed.
4. Configure environment variables, logging settings, and any other necessary configurations based on your deployment environment.
5. Start building your application by focusing on business logic implementation while leveraging the structured foundation provided by the skeleton code.

## Contributing
Contributions to this skeleton codebase are welcome! If you have suggestions for improvements, bug fixes, or new features, please feel free to submit a pull request or open an issue on GitHub.

## License
This skeleton code is open source and available under the MIT License.
