import { z } from 'zod';

export const tradeSchema = z.object({
  shares: z.coerce.number().int().positive('Enter a valid number of shares'),
});

export type TradeFormValues = z.infer<typeof tradeSchema>;

export const orderSchema = z
  .object({
    ticker: z.string().min(1, 'Ticker is required').max(6).toUpperCase(),
    side: z.enum(['buy', 'sell']),
    order_type: z.enum(['limit', 'stop', 'stop_limit']),
    shares: z.coerce.number().int().positive('Enter a valid number of shares'),
    limit_price: z.string().optional(),
    stop_price: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    if ((data.order_type === 'limit' || data.order_type === 'stop_limit') && !data.limit_price) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Limit price is required',
        path: ['limit_price'],
      });
    }
    if ((data.order_type === 'stop' || data.order_type === 'stop_limit') && !data.stop_price) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Stop price is required',
        path: ['stop_price'],
      });
    }
  });

export type OrderFormValues = z.infer<typeof orderSchema>;

export const alertSchema = z.object({
  ticker: z.string().min(1, 'Ticker is required').max(6).toUpperCase(),
  condition: z.enum(['above', 'below']),
  target_price: z.string().min(1, 'Target price is required').regex(/^\d+(\.\d{1,2})?$/, 'Invalid price'),
});

export type AlertFormValues = z.infer<typeof alertSchema>;

export const profileSchema = z.object({
  display_name: z.string().min(1, 'Name is required').max(30, 'Name is too long'),
});

export type ProfileFormValues = z.infer<typeof profileSchema>;
