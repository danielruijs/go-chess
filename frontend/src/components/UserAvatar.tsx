type Size = "sm" | "lg";

interface UserAvatarProps {
  name: string;
  size: Size;
}

const sizeClasses: Record<Size, string> = {
  sm: "w-8 h-8 rounded-lg text-sm font-bold shadow-2xs",
  lg: "w-20 h-20 rounded-2xl text-3xl font-extrabold shadow-md",
};

export function UserAvatar({ name, size }: UserAvatarProps) {
  const initial = name.trim().charAt(0).toUpperCase() || "?";

  return (
    <div
      className={`bg-gradient-to-tr from-go to-blue-600 shrink-0 flex items-center justify-center text-white uppercase select-none ${sizeClasses[size]}`}
      aria-hidden="true"
    >
      {initial}
    </div>
  );
}

export default UserAvatar;
