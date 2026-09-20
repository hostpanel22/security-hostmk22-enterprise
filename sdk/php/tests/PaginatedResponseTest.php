<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\Types\PaginatedResponse;
use PHPUnit\Framework\TestCase;

final class PaginatedResponseTest extends TestCase
{
    public function testCanCreatePaginatedResponse(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [
                ['id' => 1],
                ['id' => 2],
            ],
            1,
            10,
            25,
            'Success'
        );

        self::assertSame(
            1,
            $response->getPage()
        );

        self::assertSame(
            10,
            $response->getPerPage()
        );

        self::assertSame(
            25,
            $response->getTotal()
        );

        self::assertSame(
            3,
            $response->getTotalPages()
        );
    }

    public function testHasNextPageReturnsTrue(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [],
            1,
            10,
            30
        );

        self::assertTrue(
            $response->hasNextPage()
        );
    }

    public function testHasNextPageReturnsFalse(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [],
            3,
            10,
            30
        );

        self::assertFalse(
            $response->hasNextPage()
        );
    }

    public function testHasPreviousPageReturnsTrue(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [],
            2,
            10,
            30
        );

        self::assertTrue(
            $response->hasPreviousPage()
        );
    }

    public function testHasPreviousPageReturnsFalse(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [],
            1,
            10,
            30
        );

        self::assertFalse(
            $response->hasPreviousPage()
        );
    }

    public function testGetPagination(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [],
            2,
            10,
            45
        );

        $pagination = $response->getPagination();

        self::assertSame(
            2,
            $pagination['page']
        );

        self::assertSame(
            10,
            $pagination['perPage']
        );

        self::assertSame(
            45,
            $pagination['total']
        );

        self::assertSame(
            5,
            $pagination['totalPages']
        );
    }

    public function testToArrayContainsPagination(): void
    {
        $response = new PaginatedResponse(
            200,
            true,
            [],
            1,
            10,
            25
        );

        $array = $response->toArray();

        self::assertArrayHasKey(
            'pagination',
            $array
        );

        self::assertSame(
            1,
            $array['pagination']['page']
        );
    }

    public function testCreateFactory(): void
    {
        $response = PaginatedResponse::create(
            [
                ['id' => 1],
            ],
            1,
            10,
            50,
            'Loaded'
        );

        self::assertTrue(
            $response->isSuccess()
        );

        self::assertSame(
            'Loaded',
            $response->getMessage()
        );

        self::assertSame(
            5,
            $response->getTotalPages()
        );
    }

    public function testFromArray(): void
    {
        $response = PaginatedResponse::fromArray(
            [
                'statusCode' => 200,
                'success' => true,
                'message' => 'Loaded',
                'data' => [
                    ['id' => 1],
                ],
                'pagination' => [
                    'page' => 2,
                    'perPage' => 10,
                    'total' => 35,
                ],
            ]
        );

        self::assertSame(
            2,
            $response->getPage()
        );

        self::assertSame(
            10,
            $response->getPerPage()
        );

        self::assertSame(
            35,
            $response->getTotal()
        );

        self::assertSame(
            4,
            $response->getTotalPages()
        );
    }
}