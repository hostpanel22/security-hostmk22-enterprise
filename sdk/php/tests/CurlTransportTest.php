<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\Transport\CurlTransport;
use PHPUnit\Framework\TestCase;
use RuntimeException;

final class CurlTransportTest extends TestCase
{
    public function testDefaultConfiguration(): void
    {
        $transport = new CurlTransport();

        $this->assertTrue(
            $transport->isSslVerificationEnabled()
        );

        $this->assertSame(
            5,
            $transport->getMaxRedirects()
        );
    }

    public function testCustomConfiguration(): void
    {
        $transport = new CurlTransport(
            false,
            10
        );

        $this->assertFalse(
            $transport->isSslVerificationEnabled()
        );

        $this->assertSame(
            10,
            $transport->getMaxRedirects()
        );
    }

    public function testThrowsWhenUrlIsEmpty(): void
    {
        $transport = new CurlTransport();

        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Request URL is required.'
        );

        $transport->send(
            'GET',
            ''
        );
    }

    public function testThrowsWhenUrlContainsOnlySpaces(): void
    {
        $transport = new CurlTransport();

        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Request URL is required.'
        );

        $transport->send(
            'GET',
            '     '
        );
    }

    public function testThrowsWhenUrlIsInvalid(): void
    {
        $transport = new CurlTransport();

        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Invalid request URL.'
        );

        $transport->send(
            'GET',
            'invalid-url'
        );
    }

    public function testThrowsWhenTimeoutIsZero(): void
    {
        $transport = new CurlTransport();

        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Timeout must be greater than zero.'
        );

        $transport->send(
            'GET',
            'https://example.com',
            [],
            null,
            0
        );
    }

    public function testThrowsWhenTimeoutIsNegative(): void
    {
        $transport = new CurlTransport();

        $this->expectException(
            RuntimeException::class
        );

        $this->expectExceptionMessage(
            'Timeout must be greater than zero.'
        );

        $transport->send(
            'GET',
            'https://example.com',
            [],
            null,
            -1
        );
    }

    public function testCanInstantiateMultipleTimes(): void
    {
        $first = new CurlTransport();

        $second = new CurlTransport(
            false,
            15
        );

        $this->assertNotSame(
            $first,
            $second
        );
    }
}