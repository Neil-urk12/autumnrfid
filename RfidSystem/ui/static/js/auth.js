document.addEventListener('DOMContentLoaded', function () {
    // Handle successful login: store token
    document.body.addEventListener('loginSuccessClient', function (evt) {
        if (evt.detail && evt.detail.token) {
            localStorage.setItem('session_token', evt.detail.token);
            console.log('Token stored in localStorage. Redirecting via HX-Redirect...');
            // HX-Redirect sent by the server will handle the page navigation.
        }
    });

    // Handle logout: remove token
    document.body.addEventListener('logoutSuccessClient', function (evt) {
        localStorage.removeItem('session_token');
        console.log('Token removed from localStorage. Redirecting via HX-Redirect...');
        // HX-Redirect sent by the server will handle the page navigation.
    });

    // Function to check authentication status and redirect if necessary
    function checkAuthAndRedirect() {
        const token = localStorage.getItem('session_token');
        // Adjust path checking if your login page URL is different or has base paths
        const isLoginPage = window.location.pathname === '/login' || window.location.pathname.startsWith('/login?');

        if (token) { // User has a token
            if (isLoginPage) {
                console.log('Token found on login page, redirecting to /logs.');
                window.location.href = '/logs'; // Or use htmx.ajax for body swap
            }
            // If has token and not on login page, do nothing here.
            // The htmx:configRequest listener below will add token to requests.
        } else { // User does not have a token
            if (!isLoginPage) {
                console.log('No token found and not on login page, redirecting to login via HTMX.');
                if (window.htmx) {
                    // Ensure the URL reflects the login page to avoid confusion
                    window.history.replaceState(null, '', '/login');
                    htmx.ajax('GET', '/login', { target: 'body', swap: 'innerHTML' });
                } else {
                    console.error('HTMX not available for redirect.');
                    window.location.href = '/login'; // Fallback redirect
                }
            }
            // If no token and on login page, do nothing, let them login.
        }
    }

    // Initial check when the script loads (DOM is ready)
    checkAuthAndRedirect();

    // Configure HTMX to send the token with requests to protected endpoints
    document.body.addEventListener('htmx:configRequest', function (evt) {
        const token = localStorage.getItem('session_token');
        const requestPath = evt.detail.path;

        // Add token if it exists and the request is not for the login page itself
        // and it's a relative path (or same-origin)
        if (token && !requestPath.startsWith('/login')) {
            if (requestPath.startsWith('/') || requestPath.startsWith(window.location.origin)) {
                evt.detail.headers['Authorization'] = 'Bearer ' + token;
            }
        }
    });

    // Optional: Re-check auth after HTMX swaps content, if you navigate to full pages via HTMX
    // Be cautious with this to avoid loops, especially if login page itself could be an htmx target.
    // document.body.addEventListener('htmx:afterSwap', function(event) {
    //    // Only run check if the swap was on the body, indicating a potential page change
    //    if (event.detail.target === document.body) {
    //        checkAuthAndRedirect();
    //    }
    // });
}); 
