<?php

declare(strict_types=1);

namespace HostMK\Laravel\Exceptions;

/** API response outside the successful HTTP status range. */
final class ApiException extends HostMKException
{
    /** @param array<string,mixed>|null $responseData */
    public function __construct(
        string $message,
        private readonly int $statusCode,
        private readonly ?array $responseData = null
    ) {
        parent::__construct($message, $statusCode);
    }

    public function statusCode(): int
    {
        return $this->statusCode;
    }

    /** @return array<string,mixed>|null */
    public function responseData(): ?array
    {
        return $this->responseData;
    }
}
