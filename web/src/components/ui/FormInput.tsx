import { forwardRef } from 'react';
import { cn } from '@/lib/cn';

interface FormInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
}

export const FormInput = forwardRef<HTMLInputElement, FormInputProps>(
  ({ label, error, className, ...props }, ref) => {
    return (
      <div>
        <label className="block text-xs text-[#6E7681] mb-1.5">{label}</label>
        <input
          ref={ref}
          className={cn(
            'w-full rounded-lg bg-[#0D1117] border px-4 py-2.5 text-white placeholder-[#6E7681] focus:outline-none text-sm',
            error ? 'border-red-400 focus:border-red-400' : 'border-[#30363D] focus:border-[#50E3C2]',
            className
          )}
          {...props}
        />
        {error && <p className="text-xs text-red-400 mt-1">{error}</p>}
      </div>
    );
  }
);

FormInput.displayName = 'FormInput';
