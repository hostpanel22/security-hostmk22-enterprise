<?php

declare(strict_types=1);

namespace HostMK\Security\Types;

final class PaginatedResponse extends ApiResponse
{
    private int $page;

    private int $perPage;

    private int $total;

    private int $totalPages;

    public function __construct(
        int $statusCode,
        bool $success,
        mixed $data,
        int $page,
        int $perPage,
        int $total,
        ?string $message = null,
        array $meta = []
    ) {
        parent::__construct(
            $statusCode,
            $success,
            $data,
            $message,
            $meta
        );

        $this->page = $page;
        $this->perPage = $perPage;
        $this->total = $total;

        $this->totalPages = max(
            1,
            (int) ceil(
                $total / max(1, $perPage)
            )
        );
    }

    public function getPage(): int
    {
        return $this->page;
    }

    public function getPerPage(): int
    {
        return $this->perPage;
    }

    public function getTotal(): int
    {
        return $this->total;
    }

    public function getTotalPages(): int
    {
        return $this->totalPages;
    }

    public function hasNextPage(): bool
    {
        return $this->page < $this->totalPages;
    }

    public function hasPreviousPage(): bool
    {
        return $this->page > 1;
    }

    public function getPagination(): array
    {
        return [
            'page' => $this->page,
            'perPage' => $this->perPage,
            'total' => $this->total,
            'totalPages' => $this->totalPages,
            'hasNextPage' => $this->hasNextPage(),
            'hasPreviousPage' => $this->hasPreviousPage(),
        ];
    }

    public function toArray(): array
    {
        return array_merge(
            parent::toArray(),
            [
                'pagination' => $this->getPagination(),
            ]
        );
    }

    public static function create(
        mixed $data,
        int $page,
        int $perPage,
        int $total,
        ?string $message = null,
        int $statusCode = 200,
        array $meta = []
    ): self {
        return new self(
            $statusCode,
            true,
            $data,
            $page,
            $perPage,
            $total,
            $message,
            $meta
        );
    }

    public static function fromArray(
        array $response
    ): self {
        $pagination = $response['pagination'] ?? [];

        return new self(
            (int) ($response['statusCode'] ?? 200),
            (bool) ($response['success'] ?? true),
            $response['data'] ?? [],
            (int) ($pagination['page'] ?? 1),
            (int) ($pagination['perPage'] ?? 10),
            (int) ($pagination['total'] ?? 0),
            $response['message'] ?? null,
            (array) ($response['meta'] ?? [])
        );
    }
}