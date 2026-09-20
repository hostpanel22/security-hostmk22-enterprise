<?php

declare(strict_types=1);

namespace HostMK\Security\Exceptions;

use RuntimeException;
use Throwable;

class ApiException extends RuntimeException
{
    private ?array $responseData;

    public function __construct(
        string $message,
        int $code = 0,
        ?array $responseData = null,
        ?Throwable $previous = null
    ) {
        parent::__construct(
            $message,
            $code,
            $previous
        );

        $this->responseData = $responseData;
    }

    public function getResponseData(): ?array
    {
        return $this->responseData;
    }

    public function hasResponseData(): bool
    {
        return $this->responseData !== null;
    }

    public function getStatusCode(): int
    {
        return $this->getCode();
    }

    public function toArray(): array
    {
        return [
            'message' => $this->getMessage(),
            'statusCode' => $this->getStatusCode(),
            'responseData' => $this->responseData,
        ];
    }

    public static function badRequest(
        string $message = 'Bad Request',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            400,
            $responseData
        );
    }

    public static function unauthorized(
        string $message = 'Unauthorized',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            401,
            $responseData
        );
    }

    public static function forbidden(
        string $message = 'Forbidden',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            403,
            $responseData
        );
    }

    public static function notFound(
        string $message = 'Not Found',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            404,
            $responseData
        );
    }

    public static function conflict(
        string $message = 'Conflict',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            409,
            $responseData
        );
    }

    public static function validation(
        string $message = 'Validation Failed',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            422,
            $responseData
        );
    }

    public static function rateLimited(
        string $message = 'Too Many Requests',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            429,
            $responseData
        );
    }

    public static function serverError(
        string $message = 'Internal Server Error',
        ?array $responseData = null
    ): self {
        return new self(
            $message,
            500,
            $responseData
        );
    }
}