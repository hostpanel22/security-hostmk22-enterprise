<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\Types\ApiResponse;
use PHPUnit\Framework\TestCase;

final class ApiResponseTest extends TestCase
{
    public function testCanCreateResponse(): void
    {
        $response = new ApiResponse(
            200,
            true,
            ['id' => 1],
            'Success',
            ['requestId' => 'abc']
        );

        self::assertSame(
            200,
            $response->getStatusCode()
        );

        self::assertTrue(
            $response->isSuccess()
        );

        self::assertSame(
            ['id' => 1],
            $response->getData()
        );

        self::assertSame(
            'Success',
            $response->getMessage()
        );

        self::assertSame(
            ['requestId' => 'abc'],
            $response->getMeta()
        );
    }

    public function testHasDataReturnsTrue(): void
    {
        $response = new ApiResponse(
            200,
            true,
            ['user' => 1]
        );

        self::assertTrue(
            $response->hasData()
        );
    }

    public function testHasDataReturnsFalse(): void
    {
        $response = new ApiResponse(
            200,
            true
        );

        self::assertFalse(
            $response->hasData()
        );
    }

    public function testHasMessageReturnsTrue(): void
    {
        $response = new ApiResponse(
            200,
            true,
            null,
            'OK'
        );

        self::assertTrue(
            $response->hasMessage()
        );
    }

    public function testHasMessageReturnsFalse(): void
    {
        $response = new ApiResponse(
            200,
            true,
            null,
            ''
        );

        self::assertFalse(
            $response->hasMessage()
        );
    }

    public function testToArray(): void
    {
        $response = new ApiResponse(
            200,
            true,
            ['name' => 'john'],
            'Created',
            ['page' => 1]
        );

        $array = $response->toArray();

        self::assertSame(
            200,
            $array['statusCode']
        );

        self::assertTrue(
            $array['success']
        );

        self::assertSame(
            'Created',
            $array['message']
        );

        self::assertSame(
            ['name' => 'john'],
            $array['data']
        );

        self::assertSame(
            ['page' => 1],
            $array['meta']
        );
    }

    public function testSuccessFactory(): void
    {
        $response = ApiResponse::success(
            ['user' => 1],
            'Success',
            201
        );

        self::assertSame(
            201,
            $response->getStatusCode()
        );

        self::assertTrue(
            $response->isSuccess()
        );

        self::assertSame(
            'Success',
            $response->getMessage()
        );
    }

    public function testFailureFactory(): void
    {
        $response = ApiResponse::failure(
            'Validation failed',
            422
        );

        self::assertSame(
            422,
            $response->getStatusCode()
        );

        self::assertFalse(
            $response->isSuccess()
        );

        self::assertSame(
            'Validation failed',
            $response->getMessage()
        );
    }

    public function testFromArray(): void
    {
        $response = ApiResponse::fromArray(
            [
                'statusCode' => 201,
                'success' => true,
                'message' => 'Created',
                'data' => [
                    'id' => 10,
                ],
                'meta' => [
                    'version' => '1.0',
                ],
            ]
        );

        self::assertSame(
            201,
            $response->getStatusCode()
        );

        self::assertTrue(
            $response->isSuccess()
        );

        self::assertSame(
            'Created',
            $response->getMessage()
        );

        self::assertSame(
            ['id' => 10],
            $response->getData()
        );

        self::assertSame(
            ['version' => '1.0'],
            $response->getMeta()
        );
    }
}