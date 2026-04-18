import { Entity, PrimaryGeneratedColumn, Column, ManyToOne, OneToMany, ManyToMany, JoinTable, JoinColumn, DataSource } from 'typeorm';
import { Pharmacy1776354362570 } from '@pharmacy/src/migration/1776354362570-pharmacy';

@Entity('schedules')
export class Schedule {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ type: 'date', nullable: true })
  date: Date | null;

  @Column()
  isOpened: boolean;

  @Column({ type: 'time', nullable: true })
  timeStart: string | null;

  @Column({ type: 'time', nullable: true })
  timeEnd: string | null;
}

@Entity('producers')
export class Producer {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 40 })
  producerName: string;

  @OneToMany(() => Goods, goods => goods.producer)
  goods: Goods[];
}

@Entity('goods')
export class Goods {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 64 })
  name: string;

  @Column({ length: 64 })
  image: string;

  @Column()
  producerId: number;

  @Column()
  isInStock: boolean;

  @Column({ length: 1024, nullable: true })
  instruction: string;

  @Column({ length: 1024, nullable: true })
  description: string;

  @Column()
  isPrescriptionNeeded: boolean;

  @Column()
  price: number;

  @ManyToOne(() => Producer, producer => producer.goods)
  @JoinColumn({ name: 'producerId' })
  producer: Producer;

  @OneToMany(() => GoodsOrders, goodsOrder => goodsOrder.goods)
  goodsOrders: GoodsOrders[];

  @OneToMany(() => GoodsTag, goodsTag => goodsTag.goods)
  goodsTags: GoodsTag[];
}

@Entity('tags')
export class Tag {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 20 })
  tagName: string;

  @OneToMany(() => GoodsTag, goodsTag => goodsTag.tag)
  goodsTags: GoodsTag[];
}

@Entity('goods_tags')
export class GoodsTag {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  goodsId: number;

  @Column()
  tagId: number;

  @ManyToOne(() => Goods, goods => goods.goodsTags)
  @JoinColumn({ name: 'goodsId' })
  goods: Goods;

  @ManyToOne(() => Tag, tag => tag.goodsTags)
  @JoinColumn({ name: 'tagId' })
  tag: Tag;
}

@Entity('orders')
export class Order {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 30 })
  clientFIO: string;

  @Column({ length: 30, nullable: true })
  clientEmail: string;

  @Column({ length: 18 })
  clientPhone: string;

  @OneToMany(() => GoodsOrders, goodsOrder => goodsOrder.order, {
    onDelete: 'CASCADE'
  })
  goodsOrders: GoodsOrders[];
}

@Entity('goods_orders')
export class GoodsOrders {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  orderId: number;

  @Column()
  goodsId: number;

  @Column()
  quantity: number;

  @ManyToOne(() => Order, order => order.goodsOrders)
  @JoinColumn({ name: 'orderId' })
  order: Order;

  @ManyToOne(() => Goods, goods => goods.goodsOrders)
  @JoinColumn({ name: 'goodsId' })
  goods: Goods;
}

@Entity('announcements')
export class Announcement {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ type: 'timestamp' })
  dateTime: Date;

  @Column({ length: 64 })
  from: string;

  @Column({ length: 2048 })
  announce: string;
}

@Entity('permissions')
export class Permission {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 20 })
  action: string;

  @ManyToMany(() => Role, role => role.permissions)
  roles: Role[];
}

@Entity('roles')
export class Role {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 20 })
  name: string;

  @OneToMany(() => User, user => user.role)
  users: User[];

  @ManyToMany(() => Permission, permission => permission.roles)
  @JoinTable({
    name: 'role_permissions',
    joinColumn: { name: 'roleId', referencedColumnName: 'id' },
    inverseJoinColumn: { name: 'permissionId', referencedColumnName: 'id' }
  })
  permissions: Permission[];
}

@Entity('users')
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 40, unique: true })
  login: string;

  @Column({ length: 40 })
  userName: string;

  @Column()
  roleId: number;

  @Column({ type: 'varbinary', length: 128 })
  passwordHash: string;

  @ManyToOne(() => Role, role => role.users)
  @JoinColumn({ name: 'roleId' })
  role: Role;
}

export const AppDataSource = new DataSource({
  type: 'mysql',
  host: process.env.DB_HOST,
  port: parseInt(process.env.DB_PORT || '3306'),
  username: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
  charset: 'utf8mb4',
  synchronize: false,
  logging: false,
  entities: [ Schedule, Producer, Tag, GoodsTag, GoodsOrders, Order, Goods, Announcement, User, Role, Permission ],
  migrations: [Pharmacy1776354362570],
});