<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\Exceptions\ApiException;
use PHPUnit\Framework\TestCase;

final class ApiExceptionTest extends TestCase
{
    public function testCanCreateException(): void
    {
        $exception = new ApiException(
            'API error',
            500
        );

        $this->assertSame(
            'API error',
            $exception->getMessage()
        );

        $this->assertSame(
            500,
            $exception->getStatusCode()
        );
    }

    public function testStoresResponseData(): void
    {
        $data = [
            'error' => 'Invalid token',
        ];

        $exception = new ApiException(
            'Unauthorized',
            401,
            $data
        );

        $this->assertSame(
            $data,
            $exception->getResponseData()
        );
    }

    public function testToArray(): void
    {
        $exception = new ApiException(
            'Forbidden',
            403,
            [
                'reason' => 'missing_permission',
            ]
        );

        $result = $exception->toArray();

        $this->assertSame(
            'Forbidden',
            $result['message']
        );

        $this->assertSame(
            403,
            $result['statusCode']
        );

        $this->assertSame(
            'missing_permission',
            $result['responseData']['reason']
        );
    }

    public function testBadRequestFactory(): void
    {
        $exception = ApiException::badRequest();

        $this->assertSame(
            400,
            $exception->getStatusCode()
        );
    }

    public function testUnauthorizedFactory(): void
    {
        $exception = ApiException::unauthorized();

        $this->assertSame(
            401,
            $exception->getStatusCode()
        );
    }

    public function testForbiddenFactory(): void
    {
        $exception = ApiException::forbidden();

        $this->assertSame(
            403,
            $exception->getStatusCode()
        );
    }

    public function testNotFoundFactory(): void
    {
        $exception = ApiException::notFound();

        $this->assertSame(
            404,
            $exception->getStatusCode()
        );
    }

    public function testConflictFactory(): void
    {
        $exception = ApiException::conflict();

        $this->assertSame(
            409,
            $exception->getStatusCode()
        );
    }

    public function testValidationFactory(): void
    {
        $exception = ApiException::validation();

        $this->assertSame(
            422,
            $exception->getStatusCode()
        );
    }

    public function testRateLimitedFactory(): void
    {
        $exception = ApiException::rateLimited();

        $this->assertSame(
            429,
            $exception->getStatusCode()
        );
    }

    public function testServerErrorFactory(): void
    {
        $exception = ApiException::serverError();

        $this->assertSame(
            500,
            $exception->getStatusCode()
        );
    }
}