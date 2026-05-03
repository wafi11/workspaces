export type RoleUser = 'user' | 'admin'
export type User = {
    email : string
    name : string
    username : string
    role : RoleUser
    avatar_url : string
}