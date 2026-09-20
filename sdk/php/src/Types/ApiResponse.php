<?php

declare(strict_types=1);

namespace HostMK\Security\Types;

class ApiResponse
{
    private int $statusCode;

    private bool $success;

    private mixed $data;

    private ?string $message;

    private array $meta;

    public function __construct(
        int $statusCode,
        bool $success,
        mixed $data = null,
        ?string $message = null,
        array $meta = []
    ) {
        $this->statusCode = $statusCode;
        $this->success = $success;
        $this->data = $data;
        $this->message = $message;
        $this->meta = $meta;
    }

    public function getStatusCode(): int
    {
        return $this->statusCode;
    }

    public function isSuccess(): bool
    {
        return $this->success;
    }

    public function getData(): mixed
    {
        return $this->data;
    }

    public function getMessage(): ?string
    {
        return $this->message;
    }

    public function getMeta(): array
    {
        return $this->meta;
    }

    public function hasData(): bool
    {
        return $this->data !== null;
    }

    public function hasMessage(): bool
    {
        return $this->message !== null
            && trim($this->message) !== '';
    }

    public function toArray(): array
    {
        return [
            'statusCode' => $this->statusCode,
            'success' => $this->success,
            'message' => $this->message,
            'data' => $this->data,
            'meta' => $this->meta,
        ];
    }

    public static function success(
        mixed $data = null,
        ?string $message = null,
        int $statusCode = 200,
        array $meta = []
    ): self {
        return new self(
            $statusCode,
            true,
            $data,
            $message,
            $meta
        );
    }

    public static function failure(
        ?string $message = null,
        int $statusCode = 400,
        mixed $data = null,
        array $meta = []
    ): self {
        return new self(
            $statusCode,
            false,
            $data,
            $message,
            $meta
        );
    }

    public static function fromArray(
        array $response
    ): self {
        return new self(
            (int) ($response['statusCode'] ?? 200),
            (bool) ($response['success'] ?? true),
            $response['data'] ?? null,
            $response['message'] ?? null,
            (array) ($response['meta'] ?? [])
        );
    }
}