import { CircularProgress } from "@mui/material";

interface LoadingStateProps {
  message: string;
}

function LoadingState({ message }: LoadingStateProps) {
  return (
    <div className="flex flex-col items-center justify-center flex-1 py-20">
      <CircularProgress size={40} />
      <p className="text-sm text-slate-500 mt-4 font-medium">{message}</p>
    </div>
  );
}

export default LoadingState;
