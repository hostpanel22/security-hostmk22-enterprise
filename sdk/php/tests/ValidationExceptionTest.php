<?php

declare(strict_types=1);

namespace HostMK\Security\Tests;

use HostMK\Security\Exceptions\ValidationException;
use PHPUnit\Framework\TestCase;

final class ValidationExceptionTest extends TestCase
{
    public function testCanCreateValidationException(): void
    {
        $exception = new ValidationException(
            'Validation failed.'
        );

        $this->assertSame(
            'Validation failed.',
            $exception->getMessage()
        );

        $this->assertSame(
            422,
            $exception->getStatusCode()
        );
    }

    public function testStoresErrors(): void
    {
        $errors = [
            'email' => [
                'Email is required.',
            ],
            'password' => [
                'Password is required.',
            ],
        ];

        $exception = new ValidationException(
            'Validation failed.',
            $errors
        );

        $this->assertSame(
            $errors,
            $exception->getErrors()
        );
    }

    public function testHasErrorsReturnsTrue(): void
    {
        $exception = new ValidationException(
            'Validation failed.',
            [
                'email' => [
                    'Invalid email.',
                ],
            ]
        );

        $this->assertTrue(
            $exception->hasErrors()
        );
    }

    public function testHasErrorsReturnsFalse(): void
    {
        $exception = new ValidationException();

        $this->assertFalse(
            $exception->hasErrors()
        );
    }

    public function testHasFieldReturnsTrue(): void
    {
        $exception = new ValidationException(
            'Validation failed.',
            [
                'email' => [
                    'Invalid email.',
                ],
            ]
        );

        $this->assertTrue(
            $exception->hasField(
                'email'
            )
        );
    }

    public function testHasFieldReturnsFalse(): void
    {
        $exception = new ValidationException(
            'Validation failed.',
            [
                'email' => [
                    'Invalid email.',
                ],
            ]
        );

        $this->assertFalse(
            $exception->hasField(
                'password'
            )
        );
    }

    public function testGetFieldErrors(): void
    {
        $exception = new ValidationException(
            'Validation failed.',
            [
                'email' => [
                    'Invalid email.',
                    'Email already exists.',
                ],
            ]
        );

        $this->assertCount(
            2,
            $exception->getFieldErrors(
                'email'
            )
        );
    }

    public function testGetFieldErrorsReturnsEmptyArray(): void
    {
        $exception = new ValidationException();

        $this->assertSame(
            [],
            $exception->getFieldErrors(
                'unknown'
            )
        );
    }

    public function testFirstErrorReturnsFirstError(): void
    {
        $exception = new ValidationException(
            'Validation failed.',
            [
                'email' => [
                    'Invalid email.',
                ],
                'password' => [
                    'Password is required.',
                ],
            ]
        );

        $this->assertSame(
            'Invalid email.',
            $exception->firstError()
        );
    }

    public function testFirstErrorReturnsNull(): void
    {
        $exception = new ValidationException();

        $this->assertNull(
            $exception->firstError()
        );
    }

    public function testToArray(): void
    {
        $errors = [
            'email' => [
                'Invalid email.',
            ],
        ];

        $exception = new ValidationException(
            'Validation failed.',
            $errors
        );

        $data = $exception->toArray();

        $this->assertSame(
            'Validation failed.',
            $data['message']
        );

        $this->assertSame(
            422,
            $data['statusCode']
        );

        $this->assertSame(
            $errors,
            $data['errors']
        );
    }

    public function testFromErrorsFactory(): void
    {
        $errors = [
            'email' => [
                'Invalid email.',
            ],
        ];

        $exception =
            ValidationException::fromErrors(
                $errors
            );

        $this->assertSame(
            $errors,
            $exception->getErrors()
        );
    }

    public function testSingleFactory(): void
    {
        $exception =
            ValidationException::single(
                'email',
                'Invalid email.'
            );

        $this->assertTrue(
            $exception->hasField(
                'email'
            )
        );

        $this->assertSame(
            [
                'Invalid email.',
            ],
            $exception->getFieldErrors(
                'email'
            )
        );
    }
}