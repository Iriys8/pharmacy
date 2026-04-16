import { MigrationInterface, QueryRunner, Table, TableForeignKey, TableIndex } from "typeorm";

export class Pharmacy1776354362570 implements MigrationInterface {

      public async up(queryRunner: QueryRunner): Promise<void> {
    // 1. Таблица producers (не зависит от других)
    await queryRunner.createTable(new Table({
      name: 'producers',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'producerName', type: 'varchar', length: '40', isNullable: false }
      ]
    }), true);

    // 2. Таблица tags
    await queryRunner.createTable(new Table({
      name: 'tags',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'tagName', type: 'varchar', length: '20', isNullable: false }
      ]
    }), true);

    // 3. Таблица orders
    await queryRunner.createTable(new Table({
      name: 'orders',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'clientFIO', type: 'varchar', length: '30', isNullable: false },
        { name: 'clientEmail', type: 'varchar', length: '30', isNullable: true },
        { name: 'clientPhone', type: 'varchar', length: '18', isNullable: false }
      ]
    }), true);

    // 4. Таблица goods
    await queryRunner.createTable(new Table({
      name: 'goods',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'name', type: 'varchar', length: '64', isNullable: false },
        { name: 'image', type: 'varchar', length: '64', isNullable: false },
        { name: 'producerId', type: 'int', isNullable: false },
        { name: 'isInStock', type: 'boolean', isNullable: false },
        { name: 'instruction', type: 'varchar', length: '1024', isNullable: true },
        { name: 'description', type: 'varchar', length: '1024', isNullable: true },
        { name: 'isPrescriptionNeeded', type: 'boolean', isNullable: false },
        { name: 'price', type: 'int', isNullable: false }
      ]
    }), true);

    // 5. Таблица goods_orders (связующая)
    await queryRunner.createTable(new Table({
      name: 'goods_orders',
      columns: [
        { name: 'orderId', type: 'int', isPrimary: true },
        { name: 'goodsId', type: 'int', isPrimary: true },
        { name: 'quantity', type: 'int', isNullable: false }
      ]
    }), true);

    // 6. Таблица goods_tags (связующая для many-to-many)
    await queryRunner.createTable(new Table({
      name: 'goods_tags',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'goodsId', type: 'int', isNullable: false },
        { name: 'tagId', type: 'int', isNullable: false }
      ]
    }), true);

    // 7. Таблица announcements
    await queryRunner.createTable(new Table({
      name: 'announcements',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'dateTime', type: 'timestamp', isNullable: false },
        { name: 'from', type: 'varchar', length: '64', isNullable: false },
        { name: 'announce', type: 'varchar', length: '2048', isNullable: false }
      ]
    }), true);

    // 8. Таблица permissions
    await queryRunner.createTable(new Table({
      name: 'permissions',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'action', type: 'varchar', length: '20', isNullable: false }
      ]
    }), true);

    // 9. Таблица roles
    await queryRunner.createTable(new Table({
      name: 'roles',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'name', type: 'varchar', length: '20', isNullable: false }
      ]
    }), true);

    // 10. Таблица users
    await queryRunner.createTable(new Table({
      name: 'users',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'login', type: 'varchar', length: '40', isNullable: false },
        { name: 'userName', type: 'varchar', length: '40', isNullable: false },
        { name: 'roleId', type: 'int', isNullable: false },
        { name: 'passwordHash', type: 'varbinary', length: '128', isNullable: false }
      ]
    }), true);

    // 11. Таблица schedules
    await queryRunner.createTable(new Table({
      name: 'schedules',
      columns: [
        { name: 'id', type: 'int', isPrimary: true, isGenerated: true, generationStrategy: 'increment' },
        { name: 'date', type: 'date', isNullable: true },
        { name: 'isOpened', type: 'boolean', isNullable: false },
        { name: 'timeStart', type: 'time', isNullable: true },
        { name: 'timeEnd', type: 'time', isNullable: true }
      ]
    }), true);

    // 12. Таблица role_permissions (связующая для many-to-many)
    await queryRunner.createTable(new Table({
      name: 'role_permissions',
      columns: [
        { name: 'roleId', type: 'int', isPrimary: true },
        { name: 'permissionId', type: 'int', isPrimary: true }
      ]
    }), true);

    // ==================== СОЗДАНИЕ ИНДЕКСОВ ====================
    
    // Индекс для login в таблице users
    await queryRunner.createIndex('users', new TableIndex({
      name: 'IDX_users_login',
      columnNames: ['login'],
      isUnique: true
    }));

    // Индекс для roleId в таблице users
    await queryRunner.createIndex('users', new TableIndex({
      name: 'IDX_users_roleId',
      columnNames: ['roleId']
    }));

    // Индекс для producerId в таблице goods
    await queryRunner.createIndex('goods', new TableIndex({
      name: 'IDX_goods_producerId',
      columnNames: ['producerId']
    }));

    // Индексы для связующих таблиц
    await queryRunner.createIndex('goods_tags', new TableIndex({
      name: 'IDX_goods_tags_goodsId',
      columnNames: ['goodsId']
    }));

    await queryRunner.createIndex('goods_tags', new TableIndex({
      name: 'IDX_goods_tags_tagId',
      columnNames: ['tagId']
    }));

    await queryRunner.createIndex('goods_orders', new TableIndex({
      name: 'IDX_goods_orders_orderId',
      columnNames: ['orderId']
    }));

    await queryRunner.createIndex('goods_orders', new TableIndex({
      name: 'IDX_goods_orders_goodsId',
      columnNames: ['goodsId']
    }));

    // ==================== СОЗДАНИЕ ВНЕШНИХ КЛЮЧЕЙ ====================

    // Foreign key: goods -> producers
    await queryRunner.createForeignKey('goods', new TableForeignKey({
      columnNames: ['producerId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'producers',
      onDelete: 'RESTRICT'
    }));

    // Foreign key: goods_orders -> orders
    await queryRunner.createForeignKey('goods_orders', new TableForeignKey({
      columnNames: ['orderId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'orders',
      onDelete: 'CASCADE'
    }));

    // Foreign key: goods_orders -> goods
    await queryRunner.createForeignKey('goods_orders', new TableForeignKey({
      columnNames: ['goodsId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'goods',
      onDelete: 'CASCADE'
    }));

    // Foreign key: goods_tags -> goods
    await queryRunner.createForeignKey('goods_tags', new TableForeignKey({
      columnNames: ['goodsId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'goods',
      onDelete: 'CASCADE'
    }));

    // Foreign key: goods_tags -> tags
    await queryRunner.createForeignKey('goods_tags', new TableForeignKey({
      columnNames: ['tagId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'tags',
      onDelete: 'CASCADE'
    }));

    // Foreign key: users -> roles
    await queryRunner.createForeignKey('users', new TableForeignKey({
      columnNames: ['roleId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'roles',
      onDelete: 'RESTRICT'
    }));

    // Foreign key: role_permissions -> roles
    await queryRunner.createForeignKey('role_permissions', new TableForeignKey({
      columnNames: ['roleId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'roles',
      onDelete: 'CASCADE'
    }));

    // Foreign key: role_permissions -> permissions
    await queryRunner.createForeignKey('role_permissions', new TableForeignKey({
      columnNames: ['permissionId'],
      referencedColumnNames: ['id'],
      referencedTableName: 'permissions',
      onDelete: 'CASCADE'
    }));
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    // Удаляем внешние ключи
    const tables = ['goods', 'goods_orders', 'goods_tags', 'users', 'role_permissions'];
    for (const table of tables) {
      const foreignKeys = await queryRunner.getTable(table);
      if (foreignKeys) {
        for (const fk of foreignKeys.foreignKeys) {
          await queryRunner.dropForeignKey(table, fk);
        }
      }
    }

    // Удаляем индексы
    await queryRunner.dropIndex('users', 'IDX_users_login');
    await queryRunner.dropIndex('users', 'IDX_users_roleId');
    await queryRunner.dropIndex('goods', 'IDX_goods_producerId');
    await queryRunner.dropIndex('goods_tags', 'IDX_goods_tags_goodsId');
    await queryRunner.dropIndex('goods_tags', 'IDX_goods_tags_tagId');
    await queryRunner.dropIndex('goods_orders', 'IDX_goods_orders_orderId');
    await queryRunner.dropIndex('goods_orders', 'IDX_goods_orders_goodsId');

    // Удаляем таблицы в обратном порядке
    await queryRunner.dropTable('role_permissions');
    await queryRunner.dropTable('schedules');
    await queryRunner.dropTable('users');
    await queryRunner.dropTable('roles');
    await queryRunner.dropTable('permissions');
    await queryRunner.dropTable('announcements');
    await queryRunner.dropTable('goods_tags');
    await queryRunner.dropTable('goods_orders');
    await queryRunner.dropTable('goods');
    await queryRunner.dropTable('orders');
    await queryRunner.dropTable('tags');
    await queryRunner.dropTable('producers');
  }

}
