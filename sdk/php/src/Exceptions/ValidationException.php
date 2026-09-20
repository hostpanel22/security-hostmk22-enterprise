<?php

declare(strict_types=1);

namespace HostMK\Security\Exceptions;

use Throwable;

final class ValidationException extends ApiException
{
    /**
     * @var array<string, array<int, string>>
     */
    private array $errors;

    /**
     * @param array<string, array<int, string>> $errors
     */
    public function __construct(
        string $message = 'Validation failed.',
        array $errors = [],
        ?Throwable $previous = null
    ) {
        parent::__construct(
            $message,
            422,
            [
                'errors' => $errors,
            ],
            $previous
        );

        $this->errors = $errors;
    }

    /**
     * @return array<string, array<int, string>>
     */
    public function getErrors(): array
    {
        return $this->errors;
    }

    /**
     * @return array<int, string>
     */
    public function getFieldErrors(
        string $field
    ): array {
        return $this->errors[$field] ?? [];
    }

    public function hasErrors(): bool
    {
        return $this->errors !== [];
    }

    public function hasField(
        string $field
    ): bool {
        return isset(
            $this->errors[$field]
        );
    }

    public function firstError(): ?string
    {
        foreach (
            $this->errors as $fieldErrors
        ) {
            if (
                count($fieldErrors) > 0
            ) {
                return $fieldErrors[0];
            }
        }

        return null;
    }

    public function toArray(): array
    {
        return [
            'message' => $this->getMessage(),
            'statusCode' => $this->getStatusCode(),
            'errors' => $this->errors,
        ];
    }

    /**
     * @param array<string, array<int, string>> $errors
     */
    public static function fromErrors(
        array $errors,
        string $message = 'Validation failed.'
    ): self {
        return new self(
            $message,
            $errors
        );
    }

    public static function single(
        string $field,
        string $error,
        string $message = 'Validation failed.'
    ): self {
        return new self(
            $message,
            [
                $field => [
                    $error,
                ],
            ]
        );
    }
}