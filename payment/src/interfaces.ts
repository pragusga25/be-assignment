import type { PaymentAccountType } from '@prisma/client';

export interface IDepositData {
  paymentAccountId: string;
  amount: number;
  userId: string;
}

export interface IWithdrawData extends IDepositData {}

export interface ISendData {
  toPaymentAccountId: string;
  fromPaymentAccountId: string;
  amount: number;
  userId: string;
}

export interface DepositRequest {
  paymentAccountId: string;
  amount: number;
}

export interface WithdrawRequest {
  paymentAccountId: string;
  amount: number;
}

export interface SendRequest {
  fromAccountId: string;
  toAccountId: string;
  amount: number;
}

export interface ICreatePaymentAccount {
  paymentAccountId: string;
  paymentAccountType: PaymentAccountType;
  userId: string;
  balance: number;
}
