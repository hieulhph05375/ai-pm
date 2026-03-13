// frontend/src/lib/utils/permission.ts

/**
 * Decodes a JWT token without verifying the signature (safe for UI state checks)
 */
export function decodeJwt(token: string): any {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
        return JSON.parse(jsonPayload);
    } catch (e) {
        return null;
    }
}

/**
 * Checks if the user has the specified permission.
 * Admin users automatically bypass all rules.
 * 
 * Usage in Svelte component:
 * {#if hasPermission($authStore.user, $authStore.token, 'project:create')}
 *    <button>Create</button>
 * {/if}
 */
export function hasPermission(user: any, token: string | null, requiredPerm: string): boolean {
    if (!user || !token) return false;

    // Super bypass for Admin
    if (user.isAdmin) return true;

    // Parse token claims
    const claims = decodeJwt(token);
    if (!claims || !Array.isArray(claims.perms)) return false;

    return claims.perms.includes(requiredPerm);
}

/**
 * Checks if the user has a specific permission within a specific project.
 */
export function hasProjectPermission(
    user: any,
    token: string | null,
    projectMembers: any[],
    projectRoles: any[],
    requiredPerm: string
): boolean {
    if (!user || !token) return false;
    if (user.isAdmin) return true;

    // Find user's residence in this project
    const member = projectMembers.find(m => m.user_id === user.id);
    if (!member) return false;

    // Get the permissions for this member's role
    const role = projectRoles.find(r => r.id === member.project_role_id);
    if (!role || !role.permissions) return false;

    return role.permissions.some((p: any) => p.name === requiredPerm);
}
