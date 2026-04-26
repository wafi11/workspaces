export function LoadingScreen() {
  return (
    <div className="flex items-center justify-center w-full h-screen"
      style={{ background: 'var(--color-background-primary)' }}>
      <div className="flex flex-col items-center gap-3">
        <img
          src="http://192.168.1.10:9000/utils/logo.png"
          width={40} height={40}
          style={{ objectFit: 'contain' }}
        />
        <div className="flex gap-1">
          {[0,1,2].map(i => (
            <div key={i} style={{
              width: 5, height: 5, borderRadius: '50%',
              background: 'var(--color-text-tertiary)',
              animation: 'pulse 1.4s ease-in-out infinite',
              animationDelay: `${i * 0.2}s`
            }} />
          ))}
        </div>
      </div>
    </div>
  )
}