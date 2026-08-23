import { Button } from "@mui/material";

interface ErrorStateProps {
  message: string;
  onRetry: () => void;
}

function ErrorState({ message, onRetry }: ErrorStateProps) {
  return (
    <div className="flex flex-col items-center justify-center flex-1 py-20 px-4">
      <div className="bg-white rounded-2xl border border-rose-200 shadow-sm p-8 max-w-md w-full text-center">
        <h2 className="text-xl font-bold text-rose-600 mb-3">Error</h2>
        <p className="text-sm text-slate-600 mb-6">{message}</p>
        <Button variant="outlined" onClick={onRetry}>
          Try Again
        </Button>
      </div>
    </div>
  );
}

export default ErrorState;
